import clsx from 'clsx'

export default function SkillTag({ label, variant = 'primary' }) {
    return (
        <span
            className={clsx(
                'inline-block text-xs font-medium px-2.5 py-1 rounded-md border',
                variant === 'primary'
                    ? 'border-primary-200 text-primary-700 bg-white'
                    : 'border-gray-200 text-gray-600 bg-white'
            )}
        >
      {label}
    </span>
    )
}
