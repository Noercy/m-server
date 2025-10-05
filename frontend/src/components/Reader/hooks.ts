import { useEffect, useRef, useState } from "preact/hooks";
import Hammer from 'hammerjs'

/**
 * usePageNavigation
 * 
 * Handles the page state and enables page navigation based of user prefrences
 * 
 * @param pages -  Array of pages for the currently loaded volume/chapter
 * @param pagesPerView - Number of pages to be shown to the user, 1 or 2
 * @param direction - Reading direction left to right/right to left or vertical
 * @param separateFirstPage - Bool if the first page is seperated creating an offset
 * @returns The current page state, setter and callback functions to call next and prev page
 */
export function usePageNavigation(
    pages: string[], 
    pagesPerView: number, 
    direction: ReadDirection, 
    separateFirstPage: boolean
) {
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);

    useEffect(() => {
        setTotalPages(pages.length);
        if (currentPage > pages.length) {
            setCurrentPage(1);
        }
    }, [pages]);

    const goForward = () => {
        setCurrentPage((prev) => {
            if (separateFirstPage) {
                if (prev === 1) return 2;
                return Math.min(prev + pagesPerView, pages.length)
            }
            return Math.min(prev + pagesPerView, pages.length)
        });
    };

    const goBackwards = () => {
        setCurrentPage((prev) => {
        if (separateFirstPage) {
            if (prev === 2) return 1; // jump back to cover
            return Math.max(prev - pagesPerView, 1);
        }
        return Math.max(prev - pagesPerView, 1);
        });
    };

    const nextPage = 
        direction === "RightToLeft" ? goBackwards : goForward;
    const prevPage =
        direction === "RightToLeft" ? goForward : goBackwards;

    return { currentPage, setCurrentPage, nextPage, prevPage } 
}

/**
 * 
 * @param onPrev 
 * @param onNext 
 */
export function useKeyboardNavigation(
    onPrev: () => void, 
    onNext: () => void
) {
    useEffect(() => {
        const handleKey = (e: KeyboardEvent) => {
            if (e.key === 'ArrowLeft') onPrev();
            if (e.key === 'ArrowRight') onNext();
        };
        window.addEventListener('keydown', handleKey);
        return () => window.removeEventListener('keydown', handleKey);
    }, [onPrev, onNext]);
}

export function useScrollZoom(
    ref: preact.RefObject<HTMLElement>
) {
    const [scale, setScale] = useState(1);
    const zoomOutMax = 0.5;
    const zoomInMax = 2;
    const zoomSens = 0.001;

    useEffect(() => {
        const handleScroll = (e: WheelEvent) => {
            setScale(prev => {
                const next = prev - e.deltaY * zoomSens;
                return Math.min(Math.max(next, zoomOutMax), zoomInMax);
            });
        };

        const container = ref.current
        if (container) container.addEventListener("wheel", handleScroll)
        
        return () => container?.removeEventListener("wheel", handleScroll)
    }, [ref]);
    return scale;
}

export function useMousePanning(
    ref: preact.RefObject<HTMLElement>,
    scale: number
) {
    const [pos, setPos] = useState({ x:0, y: 0});
    const isPanning = useRef(false);
    const lastPos = useRef({ x: 0, y: 0 });

    useEffect(() => {
        const handleMouseDown = (e: MouseEvent) => {
            isPanning.current = true;
            lastPos.current = {x: e.clientX - pos.x, y: e.clientY - pos.y};
        };

        const handleMouseMove = (e: MouseEvent) => {
            if (!isPanning.current) return;
            setPos({
                x: e.clientX - lastPos.current.x,
                y: e.clientY - lastPos.current.y,
            });
        };

        const handleMouseUp = () => {
            isPanning.current = false;
        };

        const container = ref.current;
        if (container) container.addEventListener("mousedown", handleMouseDown);
        window.addEventListener("mousemove", handleMouseMove);
        window.addEventListener("mouseup", handleMouseUp);

        return () => {
            container?.removeEventListener("mousedown", handleMouseDown)
            window.removeEventListener("mousemove", handleMouseMove);
            window.removeEventListener("mouseup", handleMouseUp);
        }
    }, [ref, pos, scale]);

    return {pos}
}

/**
 * usePreloadPages
 *  
 * Preloads the two next pages and the previous page in img elements 
 * 
 * @param pages - Array of pages for the currently loaded volume/chapter
 * @param currentIndex - The currently active page number
 * @param pagesPerView - Number of pages to be shown to the user, 1 or 2
 * @returns void
 */
export function usePreloadPages(pages: string[], currentIndex: number, pagesPerView: number) {
    useEffect(() => {
        const nextIndex = currentIndex + pagesPerView;
        const prevIndex = currentIndex - pagesPerView;

        const preload = (i: number) => {
            if (i>= 0 && i < pages.length) {
                const img = new Image();
                img.src = pages[i];
            }
        };
        preload(nextIndex);
        preload(nextIndex + 1);
        preload(prevIndex);
    }, [pages, currentIndex, pagesPerView])
}

/**
 * useSwipe
 * 
 * Enables horizontal swipe gestures on touch screens for the given element.
 * Internally uses Hammer.js to detect swipe events.
 * 
 * @param ref - Refrence to the element being interacted/swiped
 * @param onSwipeLeft - Callback invoked when a left swipe is detected.
 * @param onSwipeRight - Callback invoked when a right swipe is detected.
 * @returns void
 * 
 * @remarks 
 * This hook attaches Hammer.js swipe listeners and cleans them up automatically
 * when the component unmounts or dependencies change.
 */
export function useSwipe(
    ref: preact.RefObject<HTMLElement>, 
    onSwipeLeft: () => void, 
    onSwipeRight: () => void
) {
    useEffect(() => {
        if (!ref.current) return;

        const hammer = new Hammer(ref.current);

        hammer.get("swipe").set({ direction: Hammer.DIRECTION_HORIZONTAL });

        hammer.on("swipeleft", onSwipeLeft);
        hammer.on("swiperight", onSwipeRight);

        return () => {
            hammer.off("swipeleft", onSwipeLeft);
            hammer.off("swiperight", onSwipeRight);
        };
    }, [ref, onSwipeLeft, onSwipeRight]);
}