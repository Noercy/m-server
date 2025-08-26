import { useEffect, useState } from "preact/hooks";
import Hammer from 'hammerjs'

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

/**
 * 
 * @param pages 
 * @param currentIndex 
 * @param pagesPerView 
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
 * 
 * @param pages 
 * @param pagesPerView 
 * @param direction 
 * @param separateFirstPage 
 * @returns void
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