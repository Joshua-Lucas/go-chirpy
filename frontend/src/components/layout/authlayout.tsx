import { ReactNode } from "react"

interface AuthLayoutProps {
  children: ReactNode
}

function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <main className="relative z-10 w-full max-w-440 px-gutter">
      <section
        className="bg-[#020617]/80 backdrop-blur-xl border border-slate-800 rounded-xl p-xl shadow-2xl shadow-teal-900/10"
      >
        {children}
      </section>
    </main>
  )
}
export default AuthLayout
