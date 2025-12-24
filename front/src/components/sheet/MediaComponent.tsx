import React, { useState, useEffect, useMemo } from 'react';
import { Component, StyleParams } from '../../types';
import { BaseComponent } from './BaseComponent';
import { VariableContext } from '../../utils/variableResolver';
import './MediaComponent.css';

interface MediaComponentProps {
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

export const MediaComponent: React.FC<MediaComponentProps> = ({
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
  const [isEditMode, setIsEditMode] = useState(false);
  const [currentPage, setCurrentPage] = useState(0);
  const [mediaFiles, setMediaFiles] = useState<string[]>(['']);
  const [currentFile, setCurrentFile] = useState<string>('');

  useEffect(() => {
    if (component.VarValue) {
      if (typeof component.VarValue === 'string') {
        try {
          const parsed = JSON.parse(component.VarValue);
          if (Array.isArray(parsed)) {
            setMediaFiles(parsed);
            setCurrentFile(parsed[currentPage] || '');
          } else {
            setMediaFiles([component.VarValue]);
            setCurrentFile(component.VarValue);
          }
        } catch {
          setMediaFiles([component.VarValue]);
          setCurrentFile(component.VarValue);
        }
      } else {
        try {
          const decoder = new TextDecoder();
          const decoded = decoder.decode(component.VarValue as Uint8Array);
          setMediaFiles([decoded]);
          setCurrentFile(decoded);
        } catch {
          setMediaFiles(['']);
          setCurrentFile('');
        }
      }
    } else {
      setMediaFiles(['']);
      setCurrentFile('');
    }
  }, [component.VarValue, currentPage]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const fileName = file.name;
      const newFiles = [...mediaFiles];
      newFiles[currentPage] = fileName;
      setMediaFiles(newFiles);
      setCurrentFile(fileName);
      
      const encoder = new TextEncoder();
      const varValue = encoder.encode(JSON.stringify(newFiles));
      onUpdate({
        ...component,
        VarValue: varValue,
      });
    }
  };

  const isMultivalue = mediaFiles.length > 1;

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
      <div className="media-component">
        {isEditMode ? (
          <div className="media-edit-mode">
            <div className="media-controls">
              <label>
                Количество значений:
                <input
                  type="number"
                  min="1"
                  value={mediaFiles.length}
                  onChange={(e) => {
                    const count = parseInt(e.target.value) || 1;
                    const newFiles = Array(count).fill('').map((_, i) => mediaFiles[i] ?? '');
                    setMediaFiles(newFiles);
                    setCurrentFile(newFiles[currentPage] || '');
                    const encoder = new TextEncoder()
                    const varValue = encoder.encode(JSON.stringify(newFiles));
                    onUpdate({
                      ...component,
                      VarValue: varValue,
                    });
                  }}
                />
              </label>
              <label>
                Загрузить файл:
                <input
                  type="file"
                  accept="image/*,video/*"
                  onChange={handleFileChange}
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
                  placeholder="Например: myMedia"
                />
              </label>
              <small>Используйте $VarName в других компонентах для ссылки на это значение</small>
            </div>
            <button onClick={() => setIsEditMode(false)} className="save-button">
              Сохранить
            </button>
          </div>
        ) : (
          <div className="media-normal-mode">
            <div className="media-display">
              {currentFile ? (
                <div className="media-file-info">
                  <span className="media-filename">{currentFile}</span>
                  <small>Файл: {currentFile}</small>
                </div>
              ) : (
                <div className="media-placeholder">
                  <span>Нажмите для загрузки файла</span>
                </div>
              )}
            </div>
            {isMultivalue && (
              <div className="page-selector">
                <button 
                  onClick={() => {
                    const newPage = Math.max(0, currentPage - 1);
                    setCurrentPage(newPage);
                    setCurrentFile(mediaFiles[newPage] || '');
                  }}
                  disabled={currentPage === 0}
                >
                  &lt;
                </button>
                <span>{currentPage + 1}</span>
                <button 
                  onClick={() => {
                    const newPage = Math.min(mediaFiles.length - 1, currentPage + 1);
                    setCurrentPage(newPage);
                    setCurrentFile(mediaFiles[newPage] || '');
                  }}
                  disabled={currentPage >= mediaFiles.length - 1}
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

