import React from "react"

type PrimaryButtonType = {
  type: "submit" | "reset" | "button",
  text: string,
  Icon?: React.ComponentType<{ className?: string }>;
}


function PrimaryButton({ type, text, Icon }: PrimaryButtonType) {
  return (
    <button
      className="w-full h-12 bg-primary  text-white font-public-sans font-bold text-sm tracking-wide rounded-lg flex items-center justify-center gap-2 active:scale-[0.98] transition-all duration-200"
      type={type}
    >
      {text}
      {Icon && <Icon className="w-5 h-5" />}
    </button>

  )
}

export default PrimaryButton
