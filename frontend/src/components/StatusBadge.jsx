const labels = {
    applied:   'Отправлен',
    reviewing: 'На рассмотрении',
    interview: 'Интервью',
    accepted:  'Принят',
    rejected:  'Отказ',
}

export default function StatusBadge({ status }) {
    return (
        <span style={{ border: '1px solid #999', fontSize: 12, padding: '2px 8px' }}>
            {labels[status] ?? status}
        </span>
    )
}