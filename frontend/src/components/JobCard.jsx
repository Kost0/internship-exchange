import { Link } from 'react-router-dom'

export default function JobCard({ listing }) {
    return (
        <Link to={`/listings/${listing.id}`} style={{ display: 'block', border: '1px solid #ccc', padding: 16, textDecoration: 'none', color: 'inherit', borderRadius: 8 }}>
            <div style={{ fontWeight: 'bold', marginBottom: 4 }}>{listing.title}</div>
            <div style={{ fontSize: 13, color: '#555', marginBottom: 4 }}>{listing.company?.name}</div>
            <div style={{ fontSize: 13, color: '#555', marginBottom: 8 }}>
                {listing.city && <span>{listing.city} · </span>}
                <span>{{ office: 'Офис', remote: 'Удалённо', hybrid: 'Гибрид' }[listing.format]}</span>
            </div>
            {(listing.salary_from || listing.salary_to) && (
                <div style={{ fontSize: 13, color: '#3e85dc', marginBottom: 8 }}>
                    {listing.salary_from ? `от ${listing.salary_from.toLocaleString('ru')} ₽` : ''}
                    {listing.salary_from && listing.salary_to ? ' — ' : ''}
                    {listing.salary_to ? `до ${listing.salary_to.toLocaleString('ru')} ₽` : ''}
                </div>
            )}
            <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                {listing.skills?.slice(0, 4).map(s => (
                    <span key={s.id} style={{ border: '1px solid #ccc', fontSize: 11, padding: '2px 8px' }}>{s.skill}</span>
                ))}
            </div>
        </Link>
    )
}