import logo from "@/assets/logo.png"
function AuthHeader() {
  return (
    <div className="flex flex-col items-center mb-xl">
      <div
        className="w-12 h-12 bg-surface flex items-center justify-center rounded-lg mb-md"
      >
        <img src={logo} alt="Chipry logo" className="w-12 h-12" />
      </div>
      <h1 className="font-h1 text-h2 text-white tracking-tight mb-xs">
        Chirpy
      </h1>
      <p className="font-body-md text-slate-400 opacity-75 text-center">
        Share thoughts. Stay factual.
      </p>
    </div>
  )
}

export default AuthHeader
