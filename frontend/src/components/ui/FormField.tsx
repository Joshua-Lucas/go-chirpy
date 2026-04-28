type FormFieldProps = {
  id: string;
  label: string;
  type?: string;
  placeholder?: string;
  error?: string;
};

function FormField({ id, type, placeholder, label, error }: FormFieldProps) {
  return (
    <div className="space-y-sm">
      <label
        className="font-label-caps text-label-caps text-slate-500 block"
        htmlFor={id}
      >{label}</label
      >
      <div className="relative group">
        <input
          className="w-full bg-slate-900/50 border-0 border-b border-slate-700 focus:border-primary focus:ring-0 text-white font-body-md transition-all duration-300 placeholder:text-slate-600 px-0 py-3"
          id={id}
          placeholder={placeholder}
          type={type}
        />

        {/* Decorative highlight bar on focus */}
        <div
          className="absolute bottom-0 left-0 h-1 w-0 bg-primary group-focus-within:w-full transition-all duration-500"
        ></div>
      </div>

      {error && (
        <p className="text-red-500 text-xs mt-1">
          {error}
        </p>
      )}
    </div>

  )

}

export default FormField;
