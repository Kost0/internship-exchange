import clsx from 'clsx'

const labels = {
    applied: 'Отправлен',
    reviewing: 'На рассмотрении',
    interview: 'Интервью',
    accepted: 'Принят',
    rejected: 'Отказ',
}

const styles = {
    applied: 'bg-primary-100 text-primary-700',
    reviewing: 'bg-yellow-100 text-yellow-700',
    interview: 'bg-purple-100 text-purple-700',
    accepted: 'bg-green-100 text-green-700',
    rejected: 'bg-red-100 text-red-700',
}

export default function StatusBadge({ status }) {
    return (
        <span className={clsx('text-xs font-medium px-3 py-1 rounded-full', styles[status])}>
      {labels[status]}
    </span>
    )
}