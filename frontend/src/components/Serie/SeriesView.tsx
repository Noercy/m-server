import { useEffect, useState } from "preact/hooks";
import { fetchSeriesById, type SeriesBig } from "../../api";
import { route } from "preact-router";
import CreatableSelect from "react-select/creatable"
import './SeriesView.css'

function TagInput() {


    return (
        <div class="tag-input">
            <CreatableSelect />
        </div>
    )
}

function SeriesDetailsModal({ serie, open, setOpen }: {serie: SeriesBig, open: boolean, setOpen: () => void}){
    useEffect(() => {
        if (open) {
            document.body.style.overflow = "hidden";
        } 

        return () => {
            document.body.style.overflow = "";
        }
    }, [open])

    return (
        <div class="modal-cover">
            <div class="modal-content-wrapper">
                <form class="modal-form" action="">
                    <input placeholder="Title" value={serie.Title} type="text" />
                    <textarea name="description" id="">{serie.Metadata.Description}</textarea>
                    <input type="text" placeholder="Author" />
                    <input type="text" placeholder="Artist" />
                    <TagInput />
                    <button>Submit</button>
                    <button class="close-button" onClick={setOpen}>x</button>
                </form>
            </div>
        </div>
    )
}

function SeriesDetails({ serie, setOpen}: {serie: SeriesBig, setOpen: () => void}) {
    return (
        <section class="series-details">
            <img 
                fetchPriority="high" 
                src={`/thumbnails/${serie.Cover}`} 
                alt="Cover" 
                height={240} width={160} 
                class="cover" />
            <div class="series-metadata">
                <h1>{serie.Title}</h1>
                <p>Volumes: {serie.Num_vol}</p>
                <p>Total pages: {serie.Num_images}</p>
                <p>Last updated: {serie.Created_at}</p>
                <p>Description: {serie.Metadata.Description}</p>
            </div>
            <button class="edit-button" onClick={setOpen}>Edit</button>
        </section>
    )
}

function VolumeList({ serie }: {serie: SeriesBig}) {
    return (
        <div class="volume-list">
            {serie.Volumes.map(v => (
                <div 
                    class="volume-card" 
                    key={v.ID} 
                    onClick={() => { 
                        if (serie) route(`/series/${serie.ID}/reader/${v.ID}`)
                    }}
                >
                    <img fetchPriority="high" src={`/thumbnails/${v.Cover}`} alt="" />
                    <h3>Vol. {v.Number}</h3>
                </div>
            ))}
        </div>
    )
}

export default function SeriesView({ id }: { id: string }) {
    const [series, setSeries] = useState<SeriesBig | null>(null);
    const [loading, setLoading] = useState(true);
    const [open, setOpen] = useState(false);

    useEffect(() => {
        setLoading(true);
        fetchSeriesById(Number(id))
            .then(data => setSeries(data))
            .then(() => console.log("Running useeffect"))
            .catch(console.error)
            .finally(() => setLoading(false));
    }, [id]); // maybe add series into the array if we want to 

    console.log(series)
    if (!series) return <div>Loading {loading}</div>;
    return (
        <div>
            <p>is open? {open.toString()}</p>
            <SeriesDetails 
                serie={series} 
                setOpen={() => setOpen((prev) => !prev)}
            />
            <VolumeList serie={series}/>
            {open && (
                <SeriesDetailsModal 
                    serie={series}
                    open={open} 
                    setOpen={() => setOpen((prev) => !prev)}
                />)}
        </div>
    );
}