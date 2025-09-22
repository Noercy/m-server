import { useState, useEffect, useRef } from "preact/hooks";
import './VirtualList.css'

interface VirtualSeriesListProps<T> {
  items: T[];
  itemWidth: number;   // 185
  itemHeight: number;  // fixed card height
  gap: number;         // 16px (1rem)
  renderItem: (item: T, index: number) => preact.JSX.Element;
}

export function VirtualList<T>({
  items,
  itemWidth,
  itemHeight,
  gap,
  renderItem,
}: VirtualSeriesListProps<T>) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [scrollTop, setScrollTop] = useState(0);
  const [containerWidth, setContainerWidth] = useState(window.innerWidth);

  // watch scroll
  useEffect(() => {
    const onScroll = () => setScrollTop(window.scrollY);
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  // watch resize
  useEffect(() => {
    const onResize = () => {
      setContainerWidth(containerRef.current?.offsetWidth || window.innerWidth);
      console.log("containerwidth:", containerWidth);
    };
    window.addEventListener("resize", onResize);
    onResize();
    return () => window.removeEventListener("resize", onResize);
  }, []);


  // layout math
  const itemsPerRow = Math.max(
    1,
    Math.floor((containerWidth + gap) / (itemWidth + gap))
  );
  const rowCount = Math.ceil(items.length / itemsPerRow);
  
  const rowHeight = itemHeight + gap;
  const totalHeight = rowCount * rowHeight;

  const viewportHeight = window.innerHeight;
  const startRow = Math.floor(scrollTop / rowHeight);
  const endRow = Math.min(
    rowCount - 1,
    Math.floor((scrollTop + viewportHeight) / rowHeight)
  );

  const visibleItems: { item: T; index: number; row: number; col: number }[] = [];
  for (let row = startRow; row <= endRow; row++) {
    for (let col = 0; col < itemsPerRow; col++) {
      const index = row * itemsPerRow + col;
      if (index < items.length) {
        visibleItems.push({ item: items[index], index, row, col });
      }
    }
  }

  return (
    <div
      ref={containerRef}
      className="virtuallist-container"
      style={{height: totalHeight}}
    >
      {visibleItems.map(({ item, index, row, col }) => {
        const left = col * (itemWidth + gap);
        const top = row * rowHeight;
        return (
          <div
            key={index}
            style={{
              position: "absolute",
              left,
              top,
              width: itemWidth,
              height: itemHeight,
            }}
          >
            {renderItem(item, index)}
          </div>
        );
      })}
    </div>
  );
}