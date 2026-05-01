import clsx from 'clsx'

export default function SkillTag({ label, variant = 'primary' }) {
    return (
        <span
            className={clsx(
                'inline-block text-xs font-medium px-3 py-1 rounded-full',
                variant === 'primary'
                    ? 'bg-primary-100 text-primary-700'
                    : 'bg-gray-100 text-gray-600'
            )}
        >
      {label}
    </span>
    )
}