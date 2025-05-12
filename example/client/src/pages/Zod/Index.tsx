import { useState } from 'react'
import { z } from 'zod'

// Schema definition for user data
const userValidationSchema = z.object({
  username: z.string().min(3, 'Username: 3 characters minimum required'),
  email: z.string().email('Please enter a valid email'),
  age: z.number().min(18, 'Age must be 18 or above'),
})

type User = z.infer<typeof userValidationSchema>

const ZodExample = () => {
  const [validationErrors, setValidationErrors] = useState<
    Record<string, string>
  >({})
  const [userData, setUserData] = useState<User>({
    username: '',
    email: '',
    age: 18,
  })

  const validateForm = (e: React.FormEvent) => {
    e.preventDefault()

    try {
      const validated = userValidationSchema.parse(userData)
      console.log('Valid data:', validated)
      setValidationErrors({})
    } catch (err) {
      if (err instanceof z.ZodError) {
        const errMap: Record<string, string> = {}
        for (const error of err.errors) {
          if (error.path) {
            errMap[error.path[0]] = error.message
          }
        }
        setValidationErrors(errMap)
      }
    }
  }

  const updateField = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setUserData((current) => ({
      ...current,
      [name]: name === 'age' ? Number(value) : value,
    }))
  }

  return (
    <div className="mx-auto max-w-md p-4">
      <h1 className="mb-4 font-bold text-2xl">Form Validation with Zod</h1>
      <form onSubmit={validateForm} className="space-y-4">
        <div>
          <label htmlFor="username" className="mb-1 block">
            Username:
          </label>
          <input
            id="username"
            placeholder="asd"
            type="text"
            name="username"
            value={userData.username}
            onChange={updateField}
            className="w-full rounded border p-2"
          />
          {validationErrors.username && (
            <p className="text-red-500 text-sm">{validationErrors.username}</p>
          )}

          <div>
            <label htmlFor="email" className="mb-1 block">
              Email:
            </label>
            <input
              id="email"
              type="email"
              name="email"
              value={userData.email}
              onChange={updateField}
              className="w-full rounded border p-2"
            />
            {validationErrors.email && (
              <p className="text-red-500 text-sm">{validationErrors.email}</p>
            )}
          </div>
          <div>
            <label htmlFor="age" className="mb-1 block">
              Age:
            </label>
            <input
              id="age"
              type="number"
              name="age"
              value={userData.age}
              onChange={updateField}
              className="w-full rounded border p-2"
            />
            {validationErrors.age && (
              <p className="text-red-500 text-sm">{validationErrors.age}</p>
            )}
          </div>

          <button
            type="submit"
            className="w-full rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600"
          >
            Submit
          </button>
        </div>
      </form>
    </div>
  )
}

export default ZodExample
