import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Sheet, SheetFilter, Component } from '../types';
import { apiService } from '../services/api';
import './SheetManagementPage.css';

export const SheetManagementPage: React.FC = () => {
  const [sheets, setSheets] = useState<Sheet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [templateFilter, setTemplateFilter] = useState<boolean | undefined>(undefined);
  const [sortNewest, setSortNewest] = useState<boolean | undefined>(undefined);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState<Sheet | null>(null);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<number | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    loadSheets();
  }, [templateFilter, sortNewest]);

  const loadSheets = async () => {
    setLoading(true);
    setError('');
    try {
      const allSheets = await apiService.getSheets();
      let filteredSheets = [...allSheets];

      if (templateFilter !== undefined) {
        filteredSheets = filteredSheets.filter(sheet => sheet.Template === templateFilter);
      }

      if (sortNewest !== undefined) {
        filteredSheets.sort((a, b) => {
          const timeA = a.LastWriteTime ? new Date(a.LastWriteTime).getTime() : 0;
          const timeB = b.LastWriteTime ? new Date(b.LastWriteTime).getTime() : 0;
          return sortNewest ? timeB - timeA : timeA - timeB;
        });
      } else {
        filteredSheets.sort((a, b) => a.Name.localeCompare(b.Name));
      }
      
      setSheets(filteredSheets);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateSheet = async () => {
    try {
      await apiService.createSheet();
      loadSheets();
      setShowCreateModal(false);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleUpdateSheet = async (sheet: Sheet) => {
    try {
      await apiService.updateSheet(sheet);
      loadSheets();
      setShowEditModal(null);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleDeleteSheet = async (sid: number) => {
    try {
      await apiService.deleteSheet(sid);
      loadSheets();
      setShowDeleteConfirm(null);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleCloneSheet = async (sid: number) => {
    try {
      await apiService.cloneSheet(sid);
      loadSheets();
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleExportSheet = async (sheet: Sheet) => {
    try {
      const sid = typeof sheet.Id === 'object' ? sheet.Id.Value : sheet.Id;
      const fullSheet = await apiService.getSheet(sid);
      const allComponents = await apiService.getComponents(sid);

      const exportData = {
        version: '1.0',
        sheet: {
          name: fullSheet.Name,
          description: fullSheet.Description,
          template: fullSheet.Template,
        },
        components: allComponents.map(comp => ({
          id: typeof comp.Id === 'object' ? comp.Id.Value : comp.Id,
          name: comp.Name,
          description: comp.Description,
          template: comp.Template,
          varName: comp.VarName,
          varType: comp.VarType,
          varValue: comp.VarValue,
          style: comp.Style,
        })),
        exportedAt: new Date().toISOString(),
      };
      
      const dataStr = JSON.stringify(exportData, null, 2);
      const dataBlob = new Blob([dataStr], { type: 'application/json' });
      const url = URL.createObjectURL(dataBlob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${sheet.Name}.json`;
      link.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleImportSheet = async (file: File) => {
    try {
      const text = await file.text();
      const importData = JSON.parse(text);

      let sheetInfo, componentsData;
      if (importData.version && importData.sheet) {
        sheetInfo = importData.sheet;
        componentsData = importData.components || [];
      } else {
        sheetInfo = importData;
        componentsData = [];
      }

      const newSheet = await apiService.createSheet();
      const updatedSheet: Sheet = {
        ...newSheet,
        Name: sheetInfo.name || sheetInfo.Name || 'Импортированный лист',
        Description: sheetInfo.description || sheetInfo.Description || '',
        Template: sheetInfo.template !== undefined ? sheetInfo.template : (sheetInfo.Template || false),
      };
      const savedSheet = await apiService.updateSheet(updatedSheet);
      const sheetId = typeof savedSheet.Id === 'object' ? savedSheet.Id.Value : savedSheet.Id;

      if (componentsData.length > 0) {
        const idMapping = new Map<number, number>();

        for (const compData of componentsData) {
          const newComponent = await apiService.createComponent(sheetId);
          const newCompId = typeof newComponent.Id === 'object' ? newComponent.Id.Value : newComponent.Id;
          const oldCompId = compData.id;
          idMapping.set(oldCompId, newCompId);

          const updatedComponent: Component = {
            ...newComponent,
            Name: compData.name || 'Компонент',
            Description: compData.description || '',
            Template: compData.template || false,
            VarName: compData.varName || '',
            VarType: compData.varType !== undefined ? compData.varType : 3,
            VarValue: compData.varValue || '',
            Style: compData.style || '',
          };
          
          await apiService.updateComponent(updatedComponent);
        }

        for (const compData of componentsData) {
          if (compData.varType === 0) {
            try {
              const oldChildIds = JSON.parse(compData.varValue || '[]');
              if (Array.isArray(oldChildIds)) {
                const newChildIds = oldChildIds.map((oldId: number) => idMapping.get(oldId)).filter((id: number | undefined) => id !== undefined);
                const newVarValue = JSON.stringify(newChildIds);
                
                const newCompId = idMapping.get(compData.id);
                if (newCompId) {
                  const container = await apiService.getComponent(sheetId, newCompId);
                  const updatedContainer: Component = {
                    ...container,
                    VarValue: newVarValue,
                  };
                  await apiService.updateComponent(updatedContainer);
                }
              }
            } catch {
            }
          }
        }
      }
      
      loadSheets();
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message || 'Ошибка при импорте листа');
    }
  };

  return (
    <div className="sheet-management-page">
      <div className="sheet-management-header">
        <h1>Управление листами</h1>
        <div className="header-actions">
          <button className="create-button" onClick={() => setShowCreateModal(true)}>
            Создать лист
          </button>
        </div>
      </div>

      {error && <div className="error-banner">{error}</div>}

      <div className="filters-section">
        <div className="filter-group">
          <label>Фильтр по шаблонам:</label>
          <select
            value={templateFilter === undefined ? 'all' : templateFilter.toString()}
            onChange={(e) => {
              const value = e.target.value;
              setTemplateFilter(value === 'all' ? undefined : value === 'true');
            }}
          >
            <option value="all">Все</option>
            <option value="true">Только шаблоны</option>
            <option value="false">Только не шаблоны</option>
          </select>
        </div>

        <div className="filter-group">
          <label>Сортировка:</label>
          <select
            value={sortNewest === undefined ? 'name' : sortNewest ? 'newest' : 'oldest'}
            onChange={(e) => {
              const value = e.target.value;
              setSortNewest(value === 'name' ? undefined : value === 'newest');
            }}
          >
            <option value="name">По имени</option>
            <option value="newest">Сначала новые</option>
            <option value="oldest">Сначала старые</option>
          </select>
        </div>
      </div>

      { loading ? (
          <div className="loading">Загрузка листов...</div> 
        ) : sheets.length === 0 ? (
          <div className="empty-state">
            <p>Листы не найдены. Создайте свой первый лист!</p>
          </div>
        ) : (
          <div className="sheets-grid">
            {sheets.map((sheet) => {
              const sheetId = typeof sheet.Id === 'object' ? sheet.Id.Value : sheet.Id;
              return (
                <div key={sheetId} className="sheet-card">
                  <div className="sheet-card-header">
                    <h3>{sheet.Name}</h3>
                    {sheet.Template && <span className="template-badge">Шаблон</span>}
                  </div>
                  <p className="sheet-description">{sheet.Description || 'Без описания'}</p>
                  {sheet.LastWriteTime && (
                    <p className="sheet-date">
                      Последнее изменение: {new Date(sheet.LastWriteTime).toLocaleDateString('ru-RU')}
                    </p>
                  )}
                  <div className="sheet-actions">
                    <button
                      className="action-button primary"
                      onClick={() => navigate(`/sheet/${sheetId}`)}
                    >
                      Открыть
                    </button>
                    <button
                      className="action-button"
                      onClick={() => setShowEditModal(sheet)}
                    >
                      Редактировать
                    </button>
                    <button
                      className="action-button"
                      onClick={() => handleCloneSheet(sheetId)}
                    >
                      Клонировать
                    </button>
                    <button
                      className="action-button"
                      onClick={() => handleExportSheet(sheet)}
                    >
                      Экспорт
                    </button>
                    <button
                      className="action-button danger"
                      onClick={() => setShowDeleteConfirm(sheetId)}
                    >
                      Удалить
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
      )}

      {showCreateModal && (
        <CreateSheetModal
          onClose={() => setShowCreateModal(false)}
          onCreate={handleCreateSheet}
        />
      )}

      {showEditModal && (
        <EditSheetModal
          sheet={showEditModal}
          onClose={() => setShowEditModal(null)}
          onUpdate={handleUpdateSheet}
        />
      )}

      {showDeleteConfirm !== null && (
        <DeleteConfirmModal
          sheetName={sheets.find(s => {
            const sid = typeof s.Id === 'object' ? s.Id.Value : s.Id;
            return sid === showDeleteConfirm;
          })?.Name || ''}
          onConfirm={() => handleDeleteSheet(showDeleteConfirm)}
          onCancel={() => setShowDeleteConfirm(null)}
        />
      )}

      <div className="import-section">
        <label className="import-label">
          Импорт листа (JSON):
          <input
            type="file"
            accept=".json"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) handleImportSheet(file);
            }}
            style={{ display: 'none' }}
          />
        </label>
      </div>
    </div>
  );
};

interface CreateSheetModalProps {
  onClose: () => void;
  onCreate: () => void;
}

const CreateSheetModal: React.FC<CreateSheetModalProps> = ({ onClose, onCreate }) => {
  const [fromTemplate, setFromTemplate] = useState(false);
  const [templates, setTemplates] = useState<Sheet[]>([]);
  const [selectedTemplate, setSelectedTemplate] = useState<number | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    if (fromTemplate) {
      apiService.getSheets().then(data => {
        const templateSheets = data.filter(s => s.Template);
        setTemplates(templateSheets);
      }).catch(() => {});
    }
  }, [fromTemplate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      if (fromTemplate && selectedTemplate !== null) {
        await apiService.cloneSheet(selectedTemplate);
        onCreate();
        onClose();
      } else {
        onCreate();
      }
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2>Создать новый лист</h2>
        <form onSubmit={handleSubmit}>
          {error && <div className="error-message">{error}</div>}
          <div className="form-group">
            <label>
              <input
                type="checkbox"
                checked={fromTemplate}
                onChange={(e) => setFromTemplate(e.target.checked)}
              />
              Создать из шаблона
            </label>
          </div>

          {fromTemplate && (
            <div className="form-group">
              <label>Выберите шаблон</label>
              <select
                value={selectedTemplate || ''}
                onChange={(e) => setSelectedTemplate(e.target.value ? parseInt(e.target.value) : null)}
                required
              >
                <option value="">Выберите шаблон...</option>
                {templates.map((t) => {
                  const tid = typeof t.Id === 'object' ? t.Id.Value : t.Id;
                  return (
                    <option key={tid} value={tid}>
                      {t.Name}
                    </option>
                  );
                })}
              </select>
            </div>
          )}

          <div className="modal-actions">
            <button type="button" onClick={onClose} className="cancel-button">
              Отмена
            </button>
            <button type="submit" className="submit-button">
              Создать
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

interface EditSheetModalProps {
  sheet: Sheet;
  onClose: () => void;
  onUpdate: (sheet: Sheet) => void;
}

const EditSheetModal: React.FC<EditSheetModalProps> = ({ sheet, onClose, onUpdate }) => {
  const [name, setName] = useState(sheet.Name);
  const [description, setDescription] = useState(sheet.Description);
  const [template, setTemplate] = useState(sheet.Template);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const updatedSheet: Sheet = {
      ...sheet,
      Name: name,
      Description: description,
      Template: template,
    };
    onUpdate(updatedSheet);
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2>Редактировать лист</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Название</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label>Описание</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>

          <div className="form-group">
            <label>
              <input
                type="checkbox"
                checked={template}
                onChange={(e) => setTemplate(e.target.checked)}
              />
              Пометить как шаблон
            </label>
          </div>

          <div className="modal-actions">
            <button type="button" onClick={onClose} className="cancel-button">
              Отмена
            </button>
            <button type="submit" className="submit-button">
              Обновить
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

interface DeleteConfirmModalProps {
  sheetName: string;
  onConfirm: () => void;
  onCancel: () => void;
}

const DeleteConfirmModal: React.FC<DeleteConfirmModalProps> = ({ sheetName, onConfirm, onCancel }) => {
  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2>Удалить лист</h2>
        <p>Вы уверены, что хотите удалить "{sheetName}"? Это действие нельзя отменить.</p>
        <div className="modal-actions">
          <button type="button" onClick={onCancel} className="cancel-button">
            Отмена
          </button>
          <button type="button" onClick={onConfirm} className="danger-button">
            Удалить
          </button>
        </div>
      </div>
    </div>
  );
};

