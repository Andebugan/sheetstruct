import React, { useState, useEffect } from 'react';
import { visitNodes } from 'typescript';
import { Component, StyleParams } from '../../types';
import { VariableContext } from '../../utils/variableResolver';
import './BaseComponent.css';

interface BaseComponentProps {
  component: Component;
  style: StyleParams;
  onUpdate: (component: Component) => void;
  onDelete: () => void;
  onClone: () => void;
  onResize: (style: StyleParams) => Promise<void>;
  onMove: (x: number, y: number, targetContainerId?: number) => Promise<void>;
  isSelected: () => boolean;
  onSelect?: () => void;
  onDeselect?: () => void;
  children?: React.ReactNode;
  variableContext?: VariableContext;
  parentContainerId?: number;
  parentContainerPosition?: { x: number; y: number };
}

export const BaseComponent: React.FC<BaseComponentProps> = ({
  component,
  style,
  onUpdate,
  onDelete,
  onClone,
  onResize,
  onMove,
  isSelected,
  onSelect,
  onDeselect,
  children,
  parentContainerId,
  parentContainerPosition,
}) => {
  const [isCollapsed, setIsCollapsed] = useState(style.collapsed || false);
  const [showMenu, setShowMenu] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 });
  const [resizeDragOffset, setResizeDragOffset] = useState({ x: 0, y: 0 });
  const [visualPosition, setVisualPosition] = useState({ x: style.x, y: style.y });
  const [visualSize, setVisualSize] = useState({ width: style.width, height: style.height, x: style.x, y: style.y, handle: '' });
  const [isHovered, setIsHovered] = useState(false);
  const [isResizing, setIsResizing] = useState(false);

  useEffect(() => {
    if (!isDragging) {
      const absX = parentContainerPosition ? style.x + parentContainerPosition.x : style.x;
      const absY = parentContainerPosition ? style.y + parentContainerPosition.y : style.y;
      setVisualPosition({ x: absX, y: absY });
    }
  }, [style.x, style.y, isDragging, parentContainerPosition]);

  const handleMouseEnter = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    setIsHovered(true)
  }

  const handleMouseLeave = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    setIsHovered(false)
  }

  const handleMouseDown = (e: React.MouseEvent) => {
    const target = e.target as HTMLElement;
    if (target.closest('.resize-handle') || 
        target.closest('button') || 
        target.closest('input') || 
        target.closest('textarea') ||
        target.closest('select') ||
        target.closest('.component-menu')) {
      return;
    }
    
    e.preventDefault();
    e.stopPropagation();

    if (!isSelected?.()) {
      onSelect?.();
    }

    if (isSelected?.()) {
      const absX = parentContainerPosition ? style.x + parentContainerPosition.x : style.x;
      const absY = parentContainerPosition ? style.y + parentContainerPosition.y : style.y;
      
      setIsDragging(true);
      setDragOffset({ 
        x: e.clientX - absX, 
        y: e.clientY - absY 
      });
      setVisualPosition({ x: absX, y: absY });
    }
  };

  useEffect(() => {
    if (isDragging) {
      let currentX = visualPosition.x;
      let currentY = visualPosition.y;
      let targetContainerId: number | undefined = parentContainerId;
  
      const handleMouseMove = (e: MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        const newX = e.clientX - dragOffset.x;
        const newY = e.clientY - dragOffset.y;
        const boundedX = Math.max(-50, newX);
        const boundedY = Math.max(-50, newY);
        currentX = boundedX;
        currentY = boundedY;
        setVisualPosition({ x: boundedX, y: boundedY });
  
        const elementUnder = document.elementFromPoint(e.clientX, e.clientY);
        if (elementUnder) {
          const containerElement = elementUnder.closest('[data-container-id]');
          if (containerElement && containerElement.classList.contains('container-component')) {
            const containerIdAttr = containerElement.getAttribute('data-container-id');
            if (containerIdAttr) {
              targetContainerId = parseInt(containerIdAttr);
            }
          } else {
            targetContainerId = undefined;
          }
        }
      };
  
      const handleMouseUp = () => {
        onMove(currentX, currentY, targetContainerId).then(() => {
          setIsDragging(false);
        });
  
        if (isSelected()) {
          onDeselect?.();
        }
      };
  
      document.addEventListener('mousemove', handleMouseMove, { passive: false });
      document.addEventListener('mouseup', handleMouseUp, { passive: false });
  
      return () => {
        document.removeEventListener('mousemove', handleMouseMove);
        document.removeEventListener('mouseup', handleMouseUp);
      };
    }

    if (isResizing) {
      const onMouseMove = (e: MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        const deltaX = e.clientX - resizeDragOffset.x;
        const deltaY = e.clientY - resizeDragOffset.y;

        let handle = visualSize.handle;
        let width = style.width;
        let height = style.height;
        let x = style.x;
        let y = style.y;
  
        if (handle.includes('e')) {
          width = Math.max(100, style.width + deltaX);
        }
        if (handle.includes('w')) {
          const widthChange = style.width - deltaX;
          if (widthChange >= 100) {
            width = widthChange;
            x = style.x + deltaX;
          }
        }
        if (handle.includes('s')) {
          height = Math.max(50, style.height + deltaY);
        }
        if (handle.includes('n')) {
          const heightChange = style.height - deltaY;
          if (heightChange >= 50) {
            height = heightChange;
            y = style.y + deltaY;
          }
        }

        setVisualSize({ width, height, x, y, handle })
      }

      const onMouseUp = () => {
        onResize({ ...style, width: visualSize.width, height: visualSize.height, x: visualSize.x, y: visualSize.y }).then(() => {
            setIsResizing(false);
        })
      };

      document.addEventListener('mousemove', onMouseMove, { passive: false });
      document.addEventListener('mouseup', onMouseUp, { passive: false });
  
      return () => {
        document.removeEventListener('mousemove', onMouseMove);
        document.removeEventListener('mouseup', onMouseUp);
      };
    }
  });

  const handleResize = (e: React.MouseEvent, handle: string) => {
    e.preventDefault();
    e.stopPropagation();

    resizeDragOffset.x = e.clientX;
    resizeDragOffset.y = e.clientY;
    
    visualSize.x = style.x;
    visualSize.y = style.y;
    visualSize.width = style.width;
    visualSize.height = style.height;
    visualSize.handle = handle;
    setIsResizing(true);
  };

  const toggleCollapse = () => {
    const newCollapsed = !isCollapsed;
    setIsCollapsed(newCollapsed);
    onResize({ ...style, collapsed: newCollapsed });
  };

  let actualX = style.x;
  let actualY = style.y;

  let actualHeight = style.height;
  let actualWidth = style.width;
  
  if (isDragging) {
    actualX = parentContainerPosition 
      ? visualPosition.x - parentContainerPosition.x
      : visualPosition.x;
    actualY = parentContainerPosition 
      ? visualPosition.y - parentContainerPosition.y
      : visualPosition.y;
  }

  if (isResizing) {
    actualX = parentContainerPosition
      ? visualSize.x - parentContainerPosition.x
      : visualSize.x;
    actualY = parentContainerPosition
      ? visualSize.y - parentContainerPosition.y
      : visualSize.y;

    actualHeight = visualSize.height;
    actualWidth = visualSize.width;
  }

  return (
    <div
      className={`base-component ${isSelected?.() || isHovered ? 'selected' : ''} ${isCollapsed ? 'collapsed' : ''} ${isDragging ? 'dragging' : ''}`}
      style={{
        left: `${actualX}px`,
        top: `${actualY}px`,
        width: `${actualWidth}px`,
        height: isCollapsed ? 'auto' : `${actualHeight}px`,
      }}
      onMouseDown={handleMouseDown}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <div className="component-header">
        <div className="component-title" onDoubleClick={toggleCollapse}>
          <span>{component.Name}</span>
          {isCollapsed && (
            <span className="component-description-tooltip">{component.Description}</span>
          )}
        </div>
        <div className="component-actions">
          <button onClick={() => setShowMenu(!showMenu)} className="menu-button">⋮</button>
          {showMenu && (
            <div className="component-menu">
              <button onClick={onClone}>Клонировать</button>
              <button onClick={toggleCollapse}>
                {isCollapsed ? 'Развернуть' : 'Свернуть'}
              </button>
              <button onClick={onDelete} className="danger">Удалить</button>
            </div>
          )}
        </div>
      </div>

      {!isCollapsed && (
        <>
          <div className="component-content"
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
            >{children}</div>
          {isHovered && !isCollapsed && (
            <>
              <div className="resize-handle ne" onMouseDown={(e) => handleResize(e, 'ne')} />
              <div className="resize-handle nw" onMouseDown={(e) => handleResize(e, 'nw')} />
              <div className="resize-handle se" onMouseDown={(e) => handleResize(e, 'se')} />
              <div className="resize-handle sw" onMouseDown={(e) => handleResize(e, 'sw')} />
              <div className="resize-handle ne" onMouseDown={(e) => handleResize(e, 'ne')} />
              <div className="resize-handle nw" onMouseDown={(e) => handleResize(e, 'nw')} />
              <div className="resize-handle se" onMouseDown={(e) => handleResize(e, 'se')} />
              <div className="resize-handle sw" onMouseDown={(e) => handleResize(e, 'sw')} />
            </>
          )}
        </>
      )}
    </div>
  );
};

