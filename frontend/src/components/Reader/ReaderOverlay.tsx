import './ReaderOverlay.css'

interface ReaderOverlayProps {
    direction: "LeftToRight" | "RightToLeft";
    onToggleDirection: () => void;
    seperateFirstPage: boolean;
    onToggleFirstPage: () => void;
}

export default function ReaderOverlay({
    direction,
    onToggleDirection,
    seperateFirstPage,
    onToggleFirstPage
}: ReaderOverlayProps) {

    return (
        <div class="overlay">
            <button class="disable-select" onClick={onToggleDirection}>
                {direction}
            </button>
            <button class="disable-select" onClick={onToggleFirstPage}>
                {seperateFirstPage ? "Merge First Page" : "Separate First Page"}
            </button>
        </div>
    )
}