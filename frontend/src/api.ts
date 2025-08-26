export interface Series {
    ID: number;
    Title: string;
    Cover: string;
}

export interface Volume {
    ID: number;
    Number: number;
    Num_images: number;
    Title: string;
    Path: string;
    Cover: string;
    Created_at: string;
}

export interface DbMetadata {
    Title_romaji: string;
    Title_english: string;
    Title_native: string;
    Description: string;
    Release_date: string | null;
    Publisher: string;
    Publication: string;
    Total_vol: number;
    Total_ch: number;
    Release_status: "Releasing" | "Finished" | "Hiatus" | "Cancelled" | "";
}

export interface SeriesBig {
    ID: number;
    Title: string;
    Path: string;
    Cover: string;
    Num_vol: number;
    Num_images: number;
    Created_at: string;
    Metadata: DbMetadata;
    Genres: string[];
    Tags: string[];
    Volumes: Volume[];
}


let seriesCache: Series[] = [];

export const fetchSeries = async (): Promise<Series[]> => {
    if (seriesCache.length > 0) return seriesCache;

    const res = await fetch("/api/allseries");
    if (!res.ok) throw new Error("Failed to fetch series");
    seriesCache = await res.json();
    return seriesCache;
}

export const fetchSeriesById = async (id: number): Promise<SeriesBig> => {
    // find in cache
    //if (seriesCache) {
    //   const s = seriesCache.find((x) => x.ID === id);
    //    if (s) return s;
   //}
    const res = await fetch(`/api/series/${id}`);
    if (!res.ok) throw new Error("Series not found");
    return await res.json();
}

