import { Link } from 'react-router-dom'
import SkillTag from './SkillTag'

const formatLabels = { office: 'Офис', remote: 'Удалённо', hybrid: 'Гибрид' }

export default function JobCard({ listing }) {
    return (
        <Link
            to={`/listings/${listing.id}`}
            className="block bg-white rounded-2xl border border-primary-100 p-5 hover:shadow-md hover:border-primary-300 transition-all"
        >
            <div className="flex items-start gap-3 mb-3">
                <div className="w-10 h-10 rounded-xl bg-primary-100 flex items-center justify-center text-primary-600 font-bold text-sm flex-shrink-0">
                    {listing.company?.name?.[0] ?? '?'}
                </div>
                <div className="min-w-0">
                    <p className="text-xs text-gray-500">{listing.company?.name}</p>
                    <h3 className="font-semibold text-gray-900 leading-tight">{listing.title}</h3>
                </div>
            </div>

            <div className="flex items-center gap-2 text-xs text-gray-500 mb-3">
                {listing.city && <span>{listing.city}</span>}
                {listing.city && <span>·</span>}
                <span>{formatLabels[listing.format]}</span>
            </div>

            {(listing.salaryFrom || listing.salaryTo) && (
                <p className="text-sm font-semibold text-primary-600 mb-3">
                    {listing.salaryFrom && `от ${listing.salaryFrom.toLocaleString('ru')} ₽`}
                    {listing.salaryFrom && listing.salaryTo && ' — '}
                    {listing.salaryTo && `до ${listing.salaryTo.toLocaleString('ru')} ₽`}
                </p>
            )}

            <div className="flex flex-wrap gap-1.5">
                {listing.skills.slice(0, 4).map((s) => (
                    <SkillTag key={s.id} label={s.skill} />
                ))}
            </div>
        </Link>
    )
}