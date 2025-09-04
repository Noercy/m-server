import { useEffect, useState } from "preact/hooks";
import { fetchSeries, type Series } from "../../api";
import { route } from "preact-router";
import './SeriesList.css'

const SeriesCardSkeleton = () => {
    return (
        <div className={"card"}>
            <div className={"image-wrapper"}></div>
            <div className={"h3"}>Loading</div>
        </div>
    );
};

export default function SeriesList() {
    const [series, setSeries] = useState<Series[]>([]);

    useEffect(() => {
        fetchSeries()
        .then(setSeries)
        .catch(console.error);
    }, []);

  const skeletonCount = 300; 

  return (
    <div className="series-list" id="series-list"> 
      {series.length === 0 ? Array.from({ length: skeletonCount }).map((_, i) => ( 
        <SeriesCardSkeleton key={i} /> )) : series.map((s) => ( 
        <div className="card" onClick={() => route(`/series/${s.ID}`)} key={s.ID}> 
          <img fetchPriority="high" src={`/thumbnails/${s.Cover}`} alt="Cover"/> 
          <p class="title">{s.Title}</p> 
        </div> ))} 
    </div>
  );
}

