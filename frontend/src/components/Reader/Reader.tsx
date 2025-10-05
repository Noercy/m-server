import { useEffect, useRef, useState } from "preact/hooks";
import { useKeyboardNavigation, usePreloadPages, usePageNavigation, useSwipe, useScrollZoom, useMousePanning } from "./hooks"
import ReaderOverlay from "./ReaderOverlay";
import './Reader.css'
import type { UserSettings } from "./type";

/*
type ReaderType = "Horizontal" | "Vertical"
type DisplayMode = "Single" | "Double" | "Auto";
type ReadDirection = "LeftToRight" | "RightToLeft" | "Vertical";

interface UserSettings {
    ReaderType:         ReaderType,
    DisplayMode:        DisplayMode,
    ReadDirection:      ReadDirection,
    SeparateFirstPage:  boolean,
}
*/ 

const u1Settings: UserSettings = {
    ReaderType:         "Horizontal",
    DisplayMode:        "Double",
    ReadDirection:      "RightToLeft",
    SeparateFirstPage:  true,
};

interface ReaderProps {
    id: string;
    vId: string;
}

/**
 * Reader
 * 
 * Main reader component for viewing pages. 
 * Includes viewer settings based on user preferences,
 * which should come from a global context.
 * 
 * @param props - 
 * @returns TSX element rendering the reader
 */
export default function Reader({ id, vId }: ReaderProps) {
    // TODO move to seperate usersetting 
    const [pages, setPages] = useState<string[]>([]);
    const [direction, setDirection] = useState<"LeftToRight" | "RightToLeft">("LeftToRight");
    const [separateFirstPage, setSeparateFirstPage] = useState(false);
    const [displayMode, setDisplayMode] = useState<"Single" | "Double" | "Auto">("Double")

    // TODO move to api
    useEffect(() => {
        fetch(`/api/series/${id}/reader/${vId}`)
            .then(res => res.json())
            .then((data) => {
                console.log("API returned:", data);
                setPages(data.images)
            })
            .catch(err => {
                console.error("Failed to load pages", err);
            });

            /*TODO  stop the scrolling on touch, fix later 
        function preventScroll(e: TouchEvent) {
            e.preventDefault();
        }
        document.body.addEventListener("touchmove", preventScroll, { passive: false });
            */ 
    }, [id, vId]);

    // adjust content width to avoid scrollbar hell
    useEffect(() => {
        const mainContent = document.getElementById("main-content");
        if (mainContent) mainContent.style.width = "100vw";

        return () => {
            if (mainContent) mainContent.style.width= "80vw";
        }
    }, [])

    if (pages.length === 0) {
        return <div>Loading</div>
    }

    const pagesPerView = 2
    const { currentPage, nextPage, prevPage } = usePageNavigation(pages, pagesPerView, direction, separateFirstPage);

    useKeyboardNavigation(prevPage, nextPage);
    usePreloadPages(pages, currentPage, pagesPerView);
    
    let displayPages: readonly string[] = [];

    // Handle first page, this is kind of broken
    if (separateFirstPage && currentPage === 1) {
        console.log("Seperate first page")
        displayPages = [pages[0]]
    } else {
        const startIndex = separateFirstPage ? currentPage : currentPage - 1;
        displayPages = pages.slice(startIndex, startIndex + 2);
    }

    if (direction === "RightToLeft") {
        displayPages = [...displayPages].reverse();
    }

    // swipe func, wipe not adjusting to direction swaps, kind of broken
    const readerRef = useRef<HTMLDivElement>(null);
    const onSwipeLeft = () => (direction === "RightToLeft" ? prevPage() : nextPage());
    const onSwipeRight = () => (direction === "RightToLeft" ? nextPage() : prevPage()); 
    useSwipe(readerRef, onSwipeLeft, onSwipeRight);

    const scale = useScrollZoom(readerRef);
    const {pos} = useMousePanning(readerRef, scale);

    return(
        <div class="reader">
            <div  
              ref={readerRef}  
              class="pageContainer" 
              style={{ 
                transform: 
                  `translate(${pos.x}px, ${pos.y}px) 
                  scale(${scale})`
              }}>
                {u1Settings.ReadDirection === "LeftToRight" 
                ? [...displayPages].reverse().map((src, i) => (
                    <img fetchPriority="high" key={i} src={src} alt={`Page ${i + 1}`} />
                ))
                : displayPages.map((src, i) => (
                    <img fetchPriority="high" key={i} src={src} alt={`Page ${i + 1}`} />
                ))}
            </div>   
            
            
            <span class="pageCounter disable-select">{currentPage} / {pages.length}</span>
            <ReaderOverlay
                direction={direction}
                onToggleDirection={() => 
                    setDirection((d) => (d === "LeftToRight" ? "RightToLeft" : "LeftToRight"))
                }
                seperateFirstPage={separateFirstPage} 
                onToggleFirstPage={() => setSeparateFirstPage((prev) => !prev)}
            />
            <div class="navigateButtonLeft"></div>
            <div class="navigateButtonRight"></div>
        </div>
    )
}