import React, { useState, useEffect, useMemo } from 'react';
import { Component, StyleParams, VarType } from '../../types';
import { BaseComponent } from './BaseComponent';
import { evaluateFormula, hasVariables, VariableContext } from '../../utils/variableResolver';
import './NumberComponent.css';

interface NumberComponentProps {
  component: Component;
  style: StyleParams;
  onUpdate: (component: Component) => void;
  onDelete: () => void;
  onClone: () => void;
  onResize: (style: StyleParams) => void;
  onMove: (x: number, y: number) => void;
  selected?: boolean;
  onSelect?: () => void;
  variableContext?: VariableContext;
}

export const NumberComponent: React.FC<NumberComponentProps> = ({
  component,
  style,
  onUpdate,
  onDelete,
  onClone,
  onResize,
  onMove,
  selected,
  onSelect,
  variableContext,
}) => {
  const [values, setValues] = useState<number[]>([]);
  const [step, setStep] = useState<number>(0);
  const [isInteger, setIsInteger] = useState(true);
  const [isEditMode, setIsEditMode] = useState(false);
  const [currentPage, setCurrentPage] = useState(0);
  const [formula, setFormula] = useState<string>('');

  useEffect(() => {
    if (component.VarValue) {
      try {
        if (typeof component.VarValue === 'string') {
          const parsed = JSON.parse(component.VarValue);
          if (Array.isArray(parsed.values)) {
            setValues(parsed.values);
            setStep(parsed.step || 0);
            setIsInteger(parsed.isInteger !== false);
            setFormula(parsed.formula || '');
          } else if (typeof parsed === 'number') {
            setValues([parsed]);
          }
        }
      } catch {
        const num = parseFloat(component.VarValue as string);
        if (!isNaN(num)) {
          setValues([num]);
        }
      }
    }
  }, [component.VarValue]);

  const computedValues = useMemo(() => {
    if (!formula || !variableContext || !hasVariables(formula)) {
      return values;
    }
    
    const componentId = typeof component.Id === 'object' ? component.Id.Value : component.Id;
    const computed = evaluateFormula(formula, variableContext, componentId);
    return [computed];
  }, [formula, variableContext, component, values]);

  const updateValue = (index: number, value: number) => {
    const newValues = [...values];
    newValues[index] = value;
    setValues(newValues);
    saveComponent(newValues);
  };

  const addValue = () => {
    const newValues = [...values, 0];
    setValues(newValues);
    saveComponent(newValues);
  };

  const removeValue = (index: number) => {
    const newValues = values.filter((_, i) => i !== index);
    setValues(newValues);
    saveComponent(newValues);
  };

  const increment = (index: number) => {
    if (step > 0) {
      const newValue = isInteger
        ? Math.round(values[index] + step)
        : values[index] + step;
      updateValue(index, newValue);
    }
  };

  const decrement = (index: number) => {
    if (step > 0) {
      const newValue = isInteger
        ? Math.round(values[index] - step)
        : values[index] - step;
      updateValue(index, newValue);
    }
  };

  const saveComponent = (newValues: number[], newFormula?: string) => {
    const varValue = JSON.stringify({
      values: newValues,
      step,
      isInteger,
      formula: newFormula !== undefined ? newFormula : formula,
    });
    onUpdate({
      ...component,
      VarValue: varValue,
    });
    if (newFormula !== undefined) {
      setFormula(newFormula);
    }
  };

  const displayValues = useMemo(() => {
    if (formula && variableContext && hasVariables(formula)) {
      return computedValues;
    }
    return values.length > 1 ? [values[currentPage]] : values;
  }, [values, currentPage, formula, variableContext, computedValues]);
  
  const isMultivalue = values.length > 1 && !formula;

  return (
    <BaseComponent
      component={component}
      style={style}
      onUpdate={onUpdate}
      onDelete={onDelete}
      onClone={onClone}
      onResize={onResize}
      onMove={onMove}
      selected={selected}
      onSelect={onSelect}
      variableContext={variableContext}
    >
      <div className="number-component">
        {isEditMode ? (
          <div className="number-edit-mode">
            <div className="number-controls">
              <label>
                Количество значений:
                <input
                  type="number"
                  min="1"
                  value={values.length}
                  onChange={(e) => {
                    const count = parseInt(e.target.value) || 1;
                    const newValues = Array(count).fill(0).map((_, i) => values[i] ?? 0);
                    setValues(newValues);
                    saveComponent(newValues);
                  }}
                  disabled={!!formula}
                />
              </label>
              <label>
                Размер шага:
                <input
                  type="number"
                  value={step}
                  onChange={(e) => {
                    const newStep = parseFloat(e.target.value) || 0;
                    setStep(newStep);
                    saveComponent(values);
                  }}
                  step="0.1"
                />
              </label>
              <label>
                <input
                  type="checkbox"
                  checked={isInteger}
                  onChange={(e) => {
                    setIsInteger(e.target.checked);
                    saveComponent(values);
                  }}
                />
                Целое число
              </label>
              <label>
                Формула (например: $VarName + 10 или $VarName[0] * 2):
                <input
                  type="text"
                  value={formula}
                  onChange={(e) => {
                    setFormula(e.target.value);
                    saveComponent(values, e.target.value);
                  }}
                  placeholder="Оставьте пустым для ручного ввода"
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
                  placeholder="Например: myVar"
                />
              </label>
              <small>Используйте $VarName в других компонентах для ссылки на это значение</small>
            </div>
            <button onClick={() => setIsEditMode(false)} className="save-button">
              Сохранить
            </button>
          </div>
        ) : (
          <div className="number-normal-mode">
            {displayValues.length === 0 ? (
              <div className="number-display">0</div>
            ) : (
              <div className={`number-display-container ${isMultivalue ? 'multivalue' : ''}`}>
                {displayValues.map((value, index) => (
                  <div key={index} className="number-display-item">
                    {step > 0 && (
                      <button 
                        className="number-arrow number-arrow-up" 
                        onClick={() => increment(isMultivalue ? currentPage : index)}
                        title="Увеличить"
                      >
                        ▲
                      </button>
                    )}
                    <div className="number-display">{value}</div>
                    {step > 0 && (
                      <button 
                        className="number-arrow number-arrow-down" 
                        onClick={() => decrement(isMultivalue ? currentPage : index)}
                        title="Уменьшить"
                      >
                        ▼
                      </button>
                    )}
                  </div>
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
            {selected && (
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

