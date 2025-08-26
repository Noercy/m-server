import { useEffect, useState } from "preact/hooks";
import { fetchSeries, type Series } from "../api";
import { route } from "preact-router";
import  VirtualList  from 'react-tiny-virtual-list'
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

    const skeletonCount = 100; 
	const itemsPerRow = 6;
	const rowHeight = 300; // adjust to card height
	const rowCount = Math.ceil(series.length / itemsPerRow);

  return (
	<div className="series-list" id="series-list"> 
		{series.length === 0 ? Array.from({ length: skeletonCount }).map((_, i) => 
	( <SeriesCardSkeleton key={i} /> 
	)) : 

    <VirtualList
      width={1000}
      height={1000}                 
      itemCount={rowCount}         
      itemSize={rowHeight}         
      overscanCount={3}
      renderItem={({ index, style }) => {
        // Slice the data for this row
        const start = index * itemsPerRow;
        const end = start + itemsPerRow;
        const rowItems = series.slice(start, end);

        return (
          <div key={index} style={{ ...style, display: "flex", gap: "1rem" }}>
            {rowItems.map((s) => (
              <div
                key={s.ID}
                style={{ flex: 1, cursor: "pointer" }}
                onClick={() => route(`/series/${s.ID}`)}
              >
                <img
                  src={`/thumbnails/${s.Cover}`}
                 /* loading="lazy" */
                  alt={s.Title}
                  width="200"
                  height="300"
                  style={{ display: "block", width: "100%" }}
                />
                <p>{s.Title}</p>
              </div>
            ))}
          </div>
        );
      }}
    />}
	</div>
  );
  
}