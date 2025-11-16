import { useAuth } from '../contexts/AuthContext'

export function Home() {
  const { user, logout } = useAuth()

  return (
    <div className="min-h-screen flex items-center justify-center p-5">
      <div className="bg-white p-10 rounded-lg shadow-lg w-full max-w-2xl">
        <h1 className="mb-6 text-3xl text-gray-800">Welcome, {user?.name}!</h1>
        <div className="bg-gray-50 p-5 rounded mb-6">
          <p className="my-2.5 text-gray-600">
            <strong className="text-gray-800">Email:</strong> {user?.email}
          </p>
          <p className="my-2.5 text-gray-600">
            <strong className="text-gray-800">User ID:</strong> {user?.id}
          </p>
          <p className="my-2.5 text-gray-600">
            <strong className="text-gray-800">Member since:</strong>{' '}
            {new Date(user?.created_at || '').toLocaleDateString()}
          </p>
        </div>
        <button
          onClick={logout}
          className="px-5 py-2.5 bg-red-600 text-white border-none rounded text-base font-medium cursor-pointer transition-colors hover:bg-red-700"
        >
          Logout
        </button>
      </div>
    </div>
  )
}
