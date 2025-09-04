import { useEffect, useState } from "preact/hooks";
import { fetchSeriesById, type SeriesBig } from "../api";
import { route } from "preact-router";
import './SeriesView.css'

export default function SeriesView({ id }: { id: string }) {
    const [series, setSeries] = useState<SeriesBig | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        setLoading(true);
        fetchSeriesById(Number(id))
            .then(data => setSeries(data))
            .catch(console.error)
            .finally(() => setLoading(false));
    }, [id]);
   
    if (!series) return <div>Loading {loading}</div>;

   console.log(series)
    return (
        <div>
            <img fetchPriority="high" src={`/thumbnails/${series.Cover}`} alt="Cover" height={240} width={160} class="cover" />
            <h2>{series.Title}</h2>
            <div class="volume-list">
                {series.Volumes.map(v => (
                    <div 
                        class="volume-card" 
                        key={v.ID} 
                        onClick={() => { if (series) route(`/series/${series.ID}/reader/${v.ID}`)}}
                    >
                        <img fetchPriority="high" src={`/thumbnails/${v.Cover}`} alt="" />
                        <h3>Vol. {v.Number}</h3>
                    </div>
                ))}
            </div>
        </div>
    );
}