import React, { useState, useEffect, useMemo } from 'react';
import ReactMarkdown from 'react-markdown';
import { Component, StyleParams } from '../../types';
import { BaseComponent } from './BaseComponent';
import { resolveVariables, VariableContext } from '../../utils/variableResolver';
import './TextComponent.css';

interface TextComponentProps {
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

export const TextComponent: React.FC<TextComponentProps> = ({
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
  const [isEditing, setIsEditing] = useState(false);
  const [text, setText] = useState('');
  const [isEditMode, setIsEditMode] = useState(false);
  const [currentPage, setCurrentPage] = useState(0);
  const [texts, setTexts] = useState<string[]>(['']);

  useEffect(() => {
    if (component.VarValue) {
      if (typeof component.VarValue === 'string') {
        try {
          const parsed = JSON.parse(component.VarValue);
          if (Array.isArray(parsed)) {
            setTexts(parsed);
            setText(parsed[currentPage] || '');
          } else {
            setTexts([component.VarValue]);
            setText(component.VarValue);
          }
        } catch {
          setTexts([component.VarValue]);
          setText(component.VarValue);
        }
      } else {
        try {
          const decoder = new TextDecoder();
          const decoded = decoder.decode(component.VarValue as Uint8Array);
          setTexts([decoded]);
          setText(decoded);
        } catch {
          setTexts(['']);
          setText('');
        }
      }
    } else {
      setTexts(['']);
      setText('');
    }
  }, [component.VarValue, currentPage]);

  const handleSave = () => {
    const newTexts = [...texts];
    newTexts[currentPage] = text;
    setTexts(newTexts);
    const varValue = texts.length > 1 ? JSON.stringify(newTexts) : text;
    onUpdate({
      ...component,
      VarValue: varValue,
    });
    setIsEditing(false);
  };

  const handleCancel = () => {
    setText(texts[currentPage] || '');
    setIsEditing(false);
  };

  const isMultivalue = texts.length > 1;
  const resolvedText = useMemo(() => {
    if (!variableContext || !text) return text;
    const componentId = typeof component.Id === 'object' ? component.Id.Value : component.Id;
    return resolveVariables(text, variableContext, componentId);
  }, [text, variableContext, component]);

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
      <div className="text-component">
        {isEditMode ? (
          <div className="text-edit-mode">
            <div className="text-controls">
              <label>
                Количество значений:
                <input
                  type="number"
                  min="1"
                  value={texts.length}
                  onChange={(e) => {
                    const count = parseInt(e.target.value) || 1;
                    const newTexts = Array(count).fill('').map((_, i) => texts[i] ?? '');
                    setTexts(newTexts);
                    setText(newTexts[currentPage] || '');
                    const varValue = count > 1 ? JSON.stringify(newTexts) : newTexts[0];
                    onUpdate({
                      ...component,
                      VarValue: varValue,
                    });
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
                  placeholder="Например: myText"
                />
              </label>
              <small>Используйте $VarName в других компонентах для ссылки на это значение</small>
            </div>
            <button onClick={() => setIsEditMode(false)} className="save-button">
              Сохранить
            </button>
          </div>
        ) : isEditing ? (
          <div className="text-editor">
            <textarea
              value={text}
              onChange={(e) => setText(e.target.value)}
              placeholder="Введите текст markdown..."
              className="text-textarea"
              style={{
                width: '100%',
                height: '100%',
                minHeight: `${Math.max(200, style.height - 100)}px`,
                padding: '8px',
                boxSizing: 'border-box',
                resize: 'none',
                border: '1px solid #ddd',
                borderRadius: '4px',
                fontFamily: 'inherit',
                fontSize: '14px',
              }}
            />
            <div className="text-editor-actions">
              <button onClick={handleSave} className="save-button">
                Сохранить
              </button>
              <button onClick={handleCancel} className="cancel-button">
                Отмена
              </button>
            </div>
          </div>
        ) : (
          <div className="text-normal-mode">
            <div 
              className="text-preview scrollbox" 
              onClick={() => setIsEditing(true)}
              style={{
                height: '100%',
                overflow: 'auto',
                cursor: 'text',
              }}
            >
              <ReactMarkdown>{resolvedText || 'Нажмите для редактирования...'}</ReactMarkdown>
            </div>
            {isMultivalue && (
              <div className="page-selector">
                <button 
                  onClick={() => {
                    const newPage = Math.max(0, currentPage - 1);
                    setCurrentPage(newPage);
                    setText(texts[newPage] || '');
                  }}
                  disabled={currentPage === 0}
                >
                  &lt;
                </button>
                <span>{currentPage + 1}</span>
                <button 
                  onClick={() => {
                    const newPage = Math.min(texts.length - 1, currentPage + 1);
                    setCurrentPage(newPage);
                    setText(texts[newPage] || '');
                  }}
                  disabled={currentPage >= texts.length - 1}
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

