package session

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

// go test -run Test_Store_getSessionID
func Test_Store_getSessionID(t *testing.T) {
	t.Parallel()
	expectedID := "test-session-id"

	// fiber instance
	app := fiber.New()

	t.Run("from cookie", func(t *testing.T) {
		t.Parallel()
		// session store
		store := NewStore()
		// fiber context
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		// set cookie
		ctx.Request().Header.SetCookie(store.Config.Extractor.Key, expectedID)

		require.Equal(t, expectedID, store.getSessionID(ctx))
	})

	t.Run("from header", func(t *testing.T) {
		t.Parallel()
		// session store
		store := NewStore(Config{
			Extractor: FromHeader("session_id"),
		})
		// fiber context
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		// set header
		ctx.Request().Header.Set(store.Config.Extractor.Key, expectedID)

		require.Equal(t, expectedID, store.getSessionID(ctx))
	})

	t.Run("from url query", func(t *testing.T) {
		t.Parallel()
		// session store
		store := NewStore(Config{
			Extractor: FromQuery("session_id"),
		})
		// fiber context
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		// set url parameter
		ctx.Request().SetRequestURI(fmt.Sprintf("/path?%s=%s", store.Config.Extractor.Key, expectedID))

		require.Equal(t, expectedID, store.getSessionID(ctx))
	})
}

// go test -run Test_Store_Get
// Regression: https://github.com/gofiber/fiber/issues/1408
// Regression: https://github.com/gofiber/fiber/security/advisories/GHSA-98j2-3j3p-fw2v
func Test_Store_Get(t *testing.T) {
	// Regression: https://github.com/gofiber/fiber/security/advisories/GHSA-98j2-3j3p-fw2v
	t.Parallel()
	unexpectedID := "test-session-id"
	// fiber instance
	app := fiber.New()
	t.Run("session should be re-generated if it is invalid", func(t *testing.T) {
		t.Parallel()
		// session store
		store := NewStore()
		// fiber context
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		// set cookie
		ctx.Request().Header.SetCookie(store.Config.Extractor.Key, unexpectedID)

		acquiredSession, err := store.Get(ctx)
		require.NoError(t, err)

		require.NotEqual(t, unexpectedID, acquiredSession.ID())
	})
}

// go test -run Test_Store_DeleteSession
func Test_Store_DeleteSession(t *testing.T) {
	t.Parallel()
	// fiber instance
	app := fiber.New()
	// session store
	store := NewStore()

	// fiber context
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Create a new session
	session, err := store.Get(ctx)
	require.NoError(t, err)

	// Save the session ID
	sessionID := session.ID()

	// Delete the session
	err = store.Delete(ctx, sessionID)
	require.NoError(t, err)

	// Try to get the session again
	session, err = store.Get(ctx)
	require.NoError(t, err)

	// The session ID should be different now, because the old session was deleted
	require.NotEqual(t, sessionID, session.ID())
}

func TestStore_Get_SessionAlreadyLoaded(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Create a new context
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Mock middleware and set it in the context
	middleware := &Middleware{}
	ctx.Locals(middlewareContextKey, middleware)

	// Create a new store
	store := &Store{}

	// Call the Get method
	sess, err := store.Get(ctx)

	// Assert that the error is ErrSessionAlreadyLoadedByMiddleware
	require.Nil(t, sess)
	require.Equal(t, ErrSessionAlreadyLoadedByMiddleware, err)
}

func TestStore_Delete(t *testing.T) {
	// Create a new store
	store := NewStore()

	t.Run("delete with empty session ID", func(t *testing.T) {
		err := store.Delete(context.Background(), "")
		require.Error(t, err)
		require.Equal(t, ErrEmptySessionID, err)
	})

	t.Run("delete non-existing session", func(t *testing.T) {
		err := store.Delete(context.Background(), "non-existing-session-id")
		require.NoError(t, err)
	})
}

func Test_Store_GetByID(t *testing.T) {
	t.Parallel()
	// Create a new store
	store := NewStore()

	t.Run("empty session ID", func(t *testing.T) {
		t.Parallel()
		sess, err := store.GetByID(context.Background(), "")
		require.Error(t, err)
		require.Nil(t, sess)
		require.Equal(t, ErrEmptySessionID, err)
	})

	t.Run("non-existent session ID", func(t *testing.T) {
		t.Parallel()
		sess, err := store.GetByID(context.Background(), "non-existent-session-id")
		require.Error(t, err)
		require.Nil(t, sess)
		require.Equal(t, ErrSessionIDNotFoundInStore, err)
	})

	t.Run("valid session ID", func(t *testing.T) {
		t.Parallel()
		app := fiber.New()
		// Create a new session
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		session, err := store.Get(ctx)
		defer session.Release()
		defer app.ReleaseCtx(ctx)
		require.NoError(t, err)

		// Save the session ID
		sessionID := session.ID()

		// Save the session
		err = session.Save()
		require.NoError(t, err)

		// Retrieve the session by ID
		retrievedSession, err := store.GetByID(context.Background(), sessionID)
		require.NoError(t, err)
		require.NotNil(t, retrievedSession)
		require.Equal(t, sessionID, retrievedSession.ID())

		// Call Save on the retrieved session
		retrievedSession.Set("key", "value")
		err = retrievedSession.Save()
		require.NoError(t, err)

		// Call Other Session methods
		require.Equal(t, "value", retrievedSession.Get("key"))
		require.False(t, retrievedSession.Fresh())

		require.NoError(t, retrievedSession.Reset())
		require.NoError(t, retrievedSession.Destroy())
		require.IsType(t, []any{}, retrievedSession.Keys())
		require.NoError(t, retrievedSession.Regenerate())
		require.NotPanics(t, func() {
			retrievedSession.Release()
		})
	})
}
