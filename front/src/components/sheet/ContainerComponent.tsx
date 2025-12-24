import React from 'react';
import { Component, StyleParams } from '../../types';
import { BaseComponent } from './BaseComponent';
import { VariableContext } from '../../utils/variableResolver';
import './ContainerComponent.css';

interface ContainerComponentProps {
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

export const ContainerComponent: React.FC<ContainerComponentProps> = ({
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
  variableContext,
  parentContainerId,
  parentContainerPosition,
}) => {
  const containerId = typeof component.Id === 'object' ? component.Id.Value : component.Id;

  return (
    <BaseComponent
      component={component}
      style={style}
      onUpdate={onUpdate}
      onDelete={onDelete}
      onClone={onClone}
      onResize={onResize}
      onMove={onMove}
      isSelected={isSelected}
      onSelect={onSelect}
      onDeselect={onDeselect}
      variableContext={variableContext}
      parentContainerId={parentContainerId}
      parentContainerPosition={parentContainerPosition}
    >
      <div 
        className="container-component"
        data-container-id={containerId}
      >
        {children || <p className="empty-container">Пустой контейнер. Добавьте компоненты сюда.</p>}
      </div>
    </BaseComponent>
  );
};

