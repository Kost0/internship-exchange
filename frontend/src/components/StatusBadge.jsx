import clsx from 'clsx'

const labels = {
    applied: 'Отправлен',
    reviewing: 'На рассмотрении',
    interview: 'Интервью',
    accepted: 'Принят',
    rejected: 'Отказ',
}

const styles = {
    applied: 'border-primary-200 text-primary-700 bg-white',
    reviewing: 'border-yellow-200 text-yellow-700 bg-white',
    interview: 'border-purple-200 text-purple-700 bg-white',
    accepted: 'border-green-200 text-green-700 bg-white',
    rejected: 'border-red-200 text-red-700 bg-white',
}

export default function StatusBadge({ status }) {
    return (
        <span className={clsx('text-xs font-medium px-2.5 py-1 rounded-md border', styles[status])}>
      {labels[status]}
    </span>
    )
}
