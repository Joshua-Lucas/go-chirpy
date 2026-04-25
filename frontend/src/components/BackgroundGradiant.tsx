function BackgroundGradient() {
  return (
    <div className="absolute inset-0 z-0 overflow-hidden">
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] rounded-full bg-primary/20 blur-[120px]" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] rounded-full bg-secondary/10 blur-[120px]" />
      <div
        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-full h-full opacity-30 pointer-events-none"
        style={{
          backgroundImage: "radial-gradient(#0d9488 0.5px, transparent 0.5px)",
          backgroundSize: "24px 24px",
        }}
      />
    </div>
  )
}

export default BackgroundGradient
