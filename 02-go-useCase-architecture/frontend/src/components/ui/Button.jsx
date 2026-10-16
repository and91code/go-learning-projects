export default function Button({ children, variant = 'primary', className = '', ...props }) {
  const styles = variant === 'secondary'
    ? 'border border-slate-300 bg-white text-slate-700 hover:bg-slate-50'
    : 'bg-indigo-600 text-white hover:bg-indigo-700'
  return (
    <button className={`inline-flex items-center justify-center rounded-lg px-4 py-2 font-medium transition disabled:opacity-50 ${styles} ${className}`} {...props}>
      {children}
    </button>
  )
}
