import React, { useState, useEffect } from 'react';
import { Component, StyleParams } from '../../types';
import { VariableContext } from '../../utils/variableResolver';
import './BaseComponent.css';

interface BaseComponentProps {
  component: Component;
  style: StyleParams;
  onUpdate: (component: Component) => void;
  onDelete: () => void;
  onClone: () => void;
  onResize: (style: StyleParams) => void;
  onMove: (x: number, y: number, targetContainerId?: number) => void;
  selected?: boolean;
  onSelect?: () => void;
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
  selected,
  onSelect,
  children,
  parentContainerId,
  parentContainerPosition,
}) => {
  const [isCollapsed, setIsCollapsed] = useState(style.collapsed || false);
  const [showMenu, setShowMenu] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 });
  const [visualPosition, setVisualPosition] = useState({ x: style.x, y: style.y });
  const [isResizing, setIsResizing] = useState(false);

  useEffect(() => {
    if (!isDragging) {
      const absX = parentContainerPosition ? style.x + parentContainerPosition.x : style.x;
      const absY = parentContainerPosition ? style.y + parentContainerPosition.y : style.y;
      setVisualPosition({ x: absX, y: absY });
    }
  }, [style.x, style.y, isDragging, parentContainerPosition]);

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

    if (!selected) {
      onSelect?.();
      return;
    }
    if (selected) {
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
    if (!isDragging) return;

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
      setIsDragging(false);
      onMove(currentX, currentY, targetContainerId);
    };

    document.addEventListener('mousemove', handleMouseMove, { passive: false });
    document.addEventListener('mouseup', handleMouseUp, { passive: false });

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDragging, dragOffset.x, dragOffset.y, onMove, parentContainerId]);

  const handleResize = (e: React.MouseEvent, handle: string) => {
    e.preventDefault();
    e.stopPropagation();
    
    setIsResizing(true);
    const startX = e.clientX;
    const startY = e.clientY;
    const startWidth = style.width;
    const startHeight = style.height;
    const startLeft = style.x;
    const startTop = style.y;

    const onMouseMove = (moveEvent: MouseEvent) => {
      moveEvent.preventDefault();
      moveEvent.stopPropagation();
      const deltaX = moveEvent.clientX - startX;
      const deltaY = moveEvent.clientY - startY;

      let newWidth = startWidth;
      let newHeight = startHeight;
      let newX = startLeft;
      let newY = startTop;

      if (handle.includes('e')) {
        newWidth = Math.max(100, startWidth + deltaX);
      }
      if (handle.includes('w')) {
        const widthChange = startWidth - deltaX;
        if (widthChange >= 100) {
          newWidth = widthChange;
          newX = startLeft + deltaX;
        }
      }
      if (handle.includes('s')) {
        newHeight = Math.max(50, startHeight + deltaY);
      }
      if (handle.includes('n')) {
        const heightChange = startHeight - deltaY;
        if (heightChange >= 50) {
          newHeight = heightChange;
          newY = startTop + deltaY;
        }
      }

      onResize({ ...style, width: newWidth, height: newHeight, x: newX, y: newY });
    };

    const onMouseUp = () => {
      setIsResizing(false);
      document.removeEventListener('mousemove', onMouseMove);
      document.removeEventListener('mouseup', onMouseUp);
    };

    document.addEventListener('mousemove', onMouseMove, { passive: false });
    document.addEventListener('mouseup', onMouseUp, { passive: false });
  };

  const toggleCollapse = () => {
    const newCollapsed = !isCollapsed;
    setIsCollapsed(newCollapsed);
    onResize({ ...style, collapsed: newCollapsed });
  };

  let actualX: number;
  let actualY: number;
  
  if (isDragging) {
    actualX = parentContainerPosition 
      ? visualPosition.x - parentContainerPosition.x
      : visualPosition.x;
    actualY = parentContainerPosition 
      ? visualPosition.y - parentContainerPosition.y
      : visualPosition.y;
  } else {
    actualX = style.x;
    actualY = style.y;
  }

  return (
    <div
      className={`base-component ${selected ? 'selected' : ''} ${isCollapsed ? 'collapsed' : ''} ${isDragging ? 'dragging' : ''}`}
      style={{
        left: `${actualX}px`,
        top: `${actualY}px`,
        width: `${style.width}px`,
        height: isCollapsed ? 'auto' : `${style.height}px`,
      }}
      onMouseDown={handleMouseDown}
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
          <div className="component-content">{children}</div>
          {selected && !isCollapsed && (
            <>
              <div className="resize-handle n" onMouseDown={(e) => handleResize(e, 'n')} />
              <div className="resize-handle s" onMouseDown={(e) => handleResize(e, 's')} />
              <div className="resize-handle e" onMouseDown={(e) => handleResize(e, 'e')} />
              <div className="resize-handle w" onMouseDown={(e) => handleResize(e, 'w')} />
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

