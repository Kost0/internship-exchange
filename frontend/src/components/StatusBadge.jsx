import clsx from 'clsx'

const labels = {
    applied:   'Отправлен',
    reviewing: 'На рассмотрении',
    interview: 'Интервью',
    accepted:  'Принят',
    rejected:  'Отказ',
}

export default function StatusBadge({ status }) {
    return (
        <span className="text-xs font-medium px-2.5 py-1 rounded-md border border-gray-300 text-gray-600 bg-white">
            {labels[status] ?? status}
        </span>
    )
}