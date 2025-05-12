async function loginWithLine() {
  // Redirect to the LINE login page
  // This URL should be the one provided by your backend for LINE login
  // Use this way because avoid redirect 2 times and stuck on cors origin
  window.location.href = 'http://localhost:8080/api/v1/auth/line/login'
}

// TODO: https://partnerpens.hashnode.dev/jwt-authentication-in-react-go
export function LoginPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-100">
      <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-md">
        <h1 className="mb-6 text-center font-bold text-2xl text-gray-800">
          Login
        </h1>

        <div className="social-login space-y-4">
          <button
            type="button"
            className="w-full rounded-lg bg-green-500 px-4 py-2 text-white transition hover:bg-green-600"
            onClick={loginWithLine}
          >
            Line
          </button>
        </div>
      </div>
    </div>
  )
}
