import SeriesList from "./SeriesList";
import HomeSidebar from "./HomeSidebar"
import HomeTopbar from "./HomeTopbar";

export default function Home() {
    return(
        <div class="home-wrapper">
            <HomeTopbar />
            <HomeSidebar />
            <SeriesList />
        </div>
       
    );
}