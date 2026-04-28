import { Outlet } from "react-router"
import BackgroundGradient from "../BackgroundGradiant"


function AuthLayout() {
  return (
    <>
      <BackgroundGradient />
      <main className="relative z-10 w-full max-w-550 px-gutter">
        <section
          className="bg-[#020617]/80 backdrop-blur-xl border border-slate-800 rounded-xl p-xl shadow-2xl shadow-teal-900/10"
        >
          <Outlet />
        </section>
      </main>
    </>
  )
}
export default AuthLayout
