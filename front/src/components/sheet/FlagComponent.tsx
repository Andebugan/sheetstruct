import React, { useState, useEffect } from 'react';
import { Component, StyleParams } from '../../types';
import { BaseComponent } from './BaseComponent';
import { VariableContext } from '../../utils/variableResolver';
import './FlagComponent.css';

interface FlagComponentProps {
  component: Component;
  style: StyleParams;
  onUpdate: (component: Component) => void;
  onDelete: () => void;
  onClone: () => void;
  onResize: (style: StyleParams) => Promise<void>;
  onMove: (x: number, y: number) => Promise<void>;
  isSelected: () => boolean;
  onSelect?: () => void;
  onDeselect?: () => void;
  variableContext?: VariableContext;
}

export const FlagComponent: React.FC<FlagComponentProps> = ({
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
  variableContext,
}) => {
  const [values, setValues] = useState<boolean[]>([]);
  const [isEditMode, setIsEditMode] = useState(false);
  const [currentPage, setCurrentPage] = useState(0);

  useEffect(() => {
    if (component.VarValue) {
      try {
          const decoder = new TextDecoder();
          const varValue = decoder.decode(component.VarValue as Uint8Array)
          const parsed = JSON.parse(varValue);
          if (Array.isArray(parsed)) {
            setValues(parsed);
          }
      } catch {
        const bool = typeof component.VarValue === 'string' && component.VarValue === 'true';
        setValues([bool]);
      }
    } else {
      setValues([false]);
    }
  }, [component.VarValue]);

  const toggleValue = (index: number) => {
    const newValues = [...values];
    newValues[index] = !newValues[index];
    setValues(newValues);
    saveComponent(newValues);
  };

  const addValue = () => {
    const newValues = [...values, false];
    setValues(newValues);
    saveComponent(newValues);
  };

  const removeValue = (index: number) => {
    const newValues = values.filter((_, i) => i !== index);
    setValues(newValues);
    saveComponent(newValues);
  };

  const saveComponent = (newValues: boolean[]) => {
    const encoder = new TextEncoder();
    const varValue = encoder.encode(JSON.stringify(newValues));
    onUpdate({
      ...component,
      VarValue: varValue,
    });
  };

  const displayValues = values.length > 1 ? [values[currentPage]] : values;
  const isMultivalue = values.length > 1;

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
    >
      <div className="flag-component">
        {isEditMode ? (
          <div className="flag-edit-mode">
            <div className="flag-controls">
              <label>
                Количество значений:
                <input
                  type="number"
                  min="1"
                  value={values.length}
                  onChange={(e) => {
                    const count = parseInt(e.target.value) || 1;
                    const newValues = Array(count).fill(false).map((_, i) => values[i] ?? false);
                    setValues(newValues);
                    saveComponent(newValues);
                  }}
                />
              </label>
            </div>
            <div className="component-variable">
              <label>
                Имя переменной компонента:
                <input
                  type="text"
                  value={component.VarName || ''}
                  onChange={(e) => {
                    onUpdate({
                      ...component,
                      VarName: e.target.value,
                    });
                  }}
                  placeholder="Например: myFlag"
                />
              </label>
              <small>Используйте $VarName в других компонентах для ссылки на это значение</small>
            </div>
            <button onClick={() => setIsEditMode(false)} className="save-button">
              Сохранить
            </button>
          </div>
        ) : (
          <div className="flag-normal-mode">
            {displayValues.length === 0 ? (
              <div className="flag-square unchecked" onClick={() => toggleValue(0)} />
            ) : (
              <div className={`flag-display-container ${isMultivalue ? 'multivalue' : ''}`}>
                {displayValues.map((value, index) => (
                  <div
                    key={index}
                    className={`flag-square ${value ? 'checked' : 'unchecked'}`}
                    onClick={() => toggleValue(isMultivalue ? currentPage : index)}
                  />
                ))}
              </div>
            )}
            {isMultivalue && (
              <div className="page-selector">
                <button 
                  onClick={() => setCurrentPage(Math.max(0, currentPage - 1))}
                  disabled={currentPage === 0}
                >
                  &lt;
                </button>
                <span>{currentPage + 1}</span>
                <button 
                  onClick={() => setCurrentPage(Math.min(values.length - 1, currentPage + 1))}
                  disabled={currentPage >= values.length - 1}
                >
                  &gt;
                </button>
              </div>
            )}
            {isSelected() && (
              <button 
                onClick={() => setIsEditMode(true)} 
                className="edit-mode-button"
                title="Редактировать"
              >
                ⚙
              </button>
            )}
          </div>
        )}
      </div>
    </BaseComponent>
  );
};

