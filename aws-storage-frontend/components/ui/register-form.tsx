"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldSeparator,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { authAPI } from "@/lib/api"

export function RegisterForm({
  className,
  ...props
}: React.ComponentProps<"form">) {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)
    setError("")

    const formData = new FormData(e.currentTarget)
    const username = formData.get("username") as string
    const email = formData.get("email") as string
    const password = formData.get("password") as string
    const passwordConfirm = formData.get("passwordConfirm") as string

    if (password !== passwordConfirm) {
      setError("Passwords do not match")
      setLoading(false)
      return
    }

    try {
      await authAPI.register({ username, email, password })
      router.push("/login?registered=true")
    } catch (err) {
      let errorMessage = "Registration failed";

      if (err && typeof err === 'object' && 'error' in err) {
        errorMessage = (err as { error: string }).error; // "email is already in use"
      }
      
      setError(errorMessage);
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className={cn("flex flex-col gap-6", className)} onSubmit={handleSubmit} {...props}>
      <FieldGroup>
        <div className="flex flex-col items-center gap-1 text-center">
          <h1 className="text-2xl font-bold">Register</h1>
          <p className="text-muted-foreground text-sm text-balance">
            Fill the inputs to create your accounts
          </p>
        </div>

        {error && (
          <div className="bg-destructive/15 text-destructive text-sm p-3 rounded-md">
            {error}
          </div>
        )}

        <Field>
          <FieldLabel htmlFor="username">Username</FieldLabel>
          <Input id="username" name="username" type="text" placeholder="Berk Karadalan" required disabled={loading} />
        </Field>
        <Field>
          <FieldLabel htmlFor="email">Email</FieldLabel>
          <Input id="email" name="email" type="email" placeholder="m@example.com" required disabled={loading} />
        </Field>
        <Field>
          <FieldLabel htmlFor="password">Password</FieldLabel>
          <Input id="password" name="password" type="password" required disabled={loading} />
        </Field>
        <Field>
          <FieldLabel htmlFor="passwordConfirm">Password Confirm</FieldLabel>
          <Input id="passwordConfirm" name="passwordConfirm" type="password" required disabled={loading} />
        </Field>
        <Field>
          <Button type="submit" disabled={loading}>
            {loading ? "Creating account..." : "Register"}
          </Button>
        </Field>
        <a href="/login" className="ml-auto text-sm underline-offset-4 hover:underline">
          Already have an account ?
        </a>
      </FieldGroup>
    </form>
  )
}