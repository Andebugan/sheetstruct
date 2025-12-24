import React, { useState, useEffect, useMemo, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Sheet, Component, VarType, StyleParams } from '../types';
import { apiService } from '../services/api';
import { NumberComponent } from '../components/sheet/NumberComponent';
import { FlagComponent } from '../components/sheet/FlagComponent';
import { TextComponent } from '../components/sheet/TextComponent';
import { ContainerComponent } from '../components/sheet/ContainerComponent';
import { MediaComponent } from '../components/sheet/MediaComponent';
import { VariableContext } from '../utils/variableResolver';
import './SheetBuildingPage.css';

export const SheetBuildingPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [sheet, setSheet] = useState<Sheet | null>(null);
  const [components, setComponents] = useState<Component[]>([]);
  const componentsRef = useRef<Component[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedComponents, setSelectedComponents] = useState<Set<number>>(new Set());
  const [showComponentMenu, setShowComponentMenu] = useState(false);
  const [showSheetMenu, setShowSheetMenu] = useState(false);

  useEffect(() => {
    if (id) {
      loadSheet();
    }
  }, [id]);

  useEffect(() => {
    componentsRef.current = components;
  }, [components]);

  useEffect(() => {
    const handleMoveToContainer = (e: CustomEvent) => {
      const { componentId, containerId } = e.detail;
      const component = componentsRef.current.find(c => {
        const cId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
        return cId === componentId;
      });
      if (component) {
        handleMoveComponent(component, component.Style ? parseStyle(component.Style).x : 0, component.Style ? parseStyle(component.Style).y : 0, containerId);
      }
    };
    
    window.addEventListener('moveToContainer', handleMoveToContainer as EventListener);
    return () => {
      window.removeEventListener('moveToContainer', handleMoveToContainer as EventListener);
    };
  }, []);

  const loadSheet = async () => {
    if (!id) return;
    const sheetId = parseInt(id);
    if (isNaN(sheetId)) {
      setError('Неверный ID листа');
      setLoading(false);
      return;
    }
    setLoading(true);
    setError('');
    try {
      const sheetData = await apiService.getSheet(sheetId);
      setSheet(sheetData);
      const componentsData = await apiService.getComponents(sheetId);
      setComponents(componentsData);
      componentsRef.current = componentsData;
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    } finally {
      setLoading(false);
    }
  };

  const parseStyle = (styleStr: string): StyleParams => {
    try {
      const parsed = JSON.parse(styleStr);
      return {
        x: parsed.x || 0,
        y: parsed.y || 0,
        width: parsed.width || 200,
        height: parsed.height || 100,
        collapsed: parsed.collapsed || false,
      };
    } catch {
      return { x: 0, y: 0, width: 200, height: 100 };
    }
  };

  const encodeStyle = (style: StyleParams): string => {
    return JSON.stringify(style);
  }

  const getDefaultValue = (varType: VarType): string => {
    switch (varType) {
      case VarType.Number:
        return JSON.stringify({ values: [0], step: 0, isInteger: true });
      case VarType.Flag:
        return JSON.stringify([false]);
      case VarType.Text:
        return '';
      case VarType.Container:
        return '[]';
      case VarType.Media:
        return '';
      default:
        return '';
    }
  };

  const handleCreateComponent = async (varType: VarType) => {
    if (!id) return;
    const sheetId = parseInt(id);
    if (isNaN(sheetId)) return;
    try {
      const newComponent = await apiService.createComponent(sheetId);
      const updatedComponent: Component = {
        ...newComponent,
        Name: `${VarType[varType]} Component`,
        Description: '',
        Template: false,
        VarName: `var_${Date.now()}`,
        VarType: varType,
        VarValue: getDefaultValue(varType),
        Style: encodeStyle({ x: 100, y: 100, width: 200, height: 150 }),
      };
      const saved = await apiService.updateComponent(updatedComponent);
      setComponents([...components, saved]);
      setShowComponentMenu(false);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleUpdateComponent = async (updatedComponent: Component) => {
    try {
      const saved = await apiService.updateComponent(updatedComponent);
      setComponents(components.map(c => {
        const cId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
        const savedId = typeof saved.Id === 'object' ? saved.Id.Value : saved.Id;
        return cId === savedId ? saved : c;
      }));
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleDeleteComponent = async (cid: number) => {
    if (!id) return;
    const sheetId = parseInt(id);
    if (isNaN(sheetId)) return;
    try {
      await apiService.deleteComponent(sheetId, cid);
      setComponents(components.filter(c => {
        const componentId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
        return componentId !== cid;
      }));
      setSelectedComponents(prev => {
        const newSet = new Set(prev);
        newSet.delete(cid);
        return newSet;
      });
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleCloneComponent = async (component: Component) => {
    if (!id) return;
    const sheetId = parseInt(id);
    if (isNaN(sheetId)) return;
    try {
      const componentId = typeof component.Id === 'object' ? component.Id.Value : component.Id;
      const cloned = await apiService.cloneComponent(sheetId, componentId);
      const style = parseStyle(cloned.Style);
      style.x += 20;
      style.y += 20;
      cloned.Style = encodeStyle(style);
      const updated = await apiService.updateComponent(cloned);
      setComponents([...components, updated]);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleResizeComponent = async (component: Component, style: StyleParams) => {
    const updated = {
      ...component,
      Style: encodeStyle(style),
    };
    await handleUpdateComponent(updated);
  };

  const handleMoveComponent = async (component: Component, x: number, y: number, targetContainerId?: number) => {
    const componentId = typeof component.Id === 'object' ? component.Id.Value : component.Id;
    const currentContainer = componentsRef.current.find(c => {
      if (c.VarType !== VarType.Container) return false;
      try {
        if (c.VarValue && typeof c.VarValue === 'string') {
          const parsed = JSON.parse(c.VarValue);
          if (Array.isArray(parsed) && parsed.includes(componentId)) {
            return true;
          }
        }
      } catch {
        return false;
      }
      return false;
    });
    
    const currentContainerId = currentContainer ? (typeof currentContainer.Id === 'object' ? currentContainer.Id.Value : currentContainer.Id) : undefined;
    let finalX = x;
    let finalY = y;
    
    if (targetContainerId !== undefined && targetContainerId !== currentContainerId) {
      const targetContainer = componentsRef.current.find(c => {
        const cId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
        return cId === targetContainerId && c.VarType === VarType.Container;
      });
      if (targetContainer) {
        const containerStyle = parseStyle(targetContainer.Style);
        finalX = x - containerStyle.x;
        finalY = y - containerStyle.y;
      }
    } else if (targetContainerId === undefined && currentContainerId !== undefined && currentContainer) {
      const containerStyle = parseStyle(currentContainer.Style);
      finalX = x + containerStyle.x;
      finalY = y + containerStyle.y;
    }

    const style = parseStyle(component.Style);
    style.x = finalX;
    style.y = finalY;
    await handleResizeComponent(component, style);

    if (targetContainerId !== currentContainerId) {
      if (currentContainer) {
        try {
          if (currentContainer.VarValue && typeof currentContainer.VarValue === 'string') {
            const parsed = JSON.parse(currentContainer.VarValue);
            if (Array.isArray(parsed)) {
              const updatedIds = parsed.filter((id: number) => id !== componentId);
              await handleUpdateComponent({
                ...currentContainer,
                VarValue: JSON.stringify(updatedIds),
              });
            }
          }
        } catch {
        }
      }

      if (targetContainerId !== undefined) {
        const targetContainer = componentsRef.current.find(c => {
          const cId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
          return cId === targetContainerId && c.VarType === VarType.Container;
        });
        
        if (targetContainer) {
          let childIds: number[] = [];
          try {
            if (targetContainer.VarValue && typeof targetContainer.VarValue === 'string') {
              const parsed = JSON.parse(targetContainer.VarValue);
              if (Array.isArray(parsed)) {
                childIds = parsed;
              }
            }
          } catch {
            childIds = [];
          }
          
          if (!childIds.includes(componentId)) {
            childIds.push(componentId);
            await handleUpdateComponent({
              ...targetContainer,
              VarValue: JSON.stringify(childIds),
            });
          }
        }
      }
    }
  };

  const handleCheckSelected = (cid: number) => {
      const isSelected = selectedComponents.has(cid)
      return isSelected;
  };

  const handleSelectComponent = (cid: number) => {
    if (!selectedComponents.has(cid)) {
      selectedComponents.add(cid);
    }
  };

  const handleDeselectComponent = (cid: number) => {
    if (selectedComponents.has(cid)) {
      selectedComponents.delete(cid);
    }
  };

  const handleSaveSheet = async () => {
    if (!sheet) return;
    try {
      await apiService.updateSheet(sheet);

      setError('');
      const successMsg = document.createElement('div');
      successMsg.textContent = 'Лист успешно сохранен!';
      successMsg.style.cssText = 'position: fixed; top: 20px; right: 20px; background: #4caf50; color: white; padding: 12px 24px; border-radius: 4px; z-index: 10000; box-shadow: 0 4px 8px rgba(0,0,0,0.2);';
      document.body.appendChild(successMsg);
      setTimeout(() => {
        document.body.removeChild(successMsg);
      }, 2000);
    } catch (err) {
      const apiError = apiService.handleError(err);
      setError(apiError.message);
    }
  };

  const handleExportSheet = async () => {
    if (!sheet) return;
    try {
      const sheetId = typeof sheet.Id === 'object' ? sheet.Id.Value : sheet.Id;
      const fullSheet = await apiService.getSheet(sheetId);

      const exportData = {
        version: '1.0',
        sheet: {
          name: fullSheet.Name,
          description: fullSheet.Description,
          template: fullSheet.Template,
        },
        components: components.map(comp => ({
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

  const variableContext: VariableContext = useMemo(() => {
    if (!sheet || !id) {
      return { components: [], sheetId: 0 };
    }
    const sheetId = typeof sheet.Id === 'object' ? sheet.Id.Value : sheet.Id;
    return {
      components,
      sheetId,
    };
  }, [components, sheet, id]);

  const renderComponent = (component: Component, renderedIds: Set<number> = new Set(), parentContainerId?: number, parentContainerPosition?: { x: number; y: number }): React.ReactNode => {
    const style = parseStyle(component.Style);
    const componentId = typeof component.Id === 'object' ? component.Id.Value : component.Id;

    if (renderedIds.has(componentId)) {
      return null;
    }
    
    const newRenderedIds = new Set(renderedIds);
    newRenderedIds.add(componentId);
    
    const commonProps = {
      component,
      style,
      onUpdate: handleUpdateComponent,
      onDelete: () => handleDeleteComponent(componentId),
      onClone: () => handleCloneComponent(component),
      onResize: (newStyle: StyleParams) => handleResizeComponent(component, newStyle),
      onMove: (x: number, y: number, targetContainerId?: number) => handleMoveComponent(component, x, y, targetContainerId),
      isSelected: () => handleCheckSelected(componentId),
      onSelect: () => handleSelectComponent(componentId),
      onDeselect: () => handleDeselectComponent(componentId),
      variableContext,
      parentContainerId,
      parentContainerPosition,
    };

    switch (component.VarType) {
      case VarType.Number:
        return <NumberComponent key={componentId} {...commonProps} />;
      case VarType.Flag:
        return <FlagComponent key={componentId} {...commonProps} />;
      case VarType.Text:
        return <TextComponent key={componentId} {...commonProps} />;
      case VarType.Media:
        return <MediaComponent key={componentId} {...commonProps} />;
      case VarType.Container:
        let childIds: number[] = [];
        try {
          if (component.VarValue && typeof component.VarValue === 'string') {
            const parsed = JSON.parse(component.VarValue);
            if (Array.isArray(parsed)) {
              childIds = parsed;
            }
          }
        } catch {
        }
        
        const containerPosition = { x: style.x, y: style.y };
        
        return (
          <ContainerComponent key={componentId} {...commonProps}>
            {childIds.length > 0 ? (
              childIds.map(childId => {
                const childComponent = components.find(c => {
                  const cId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
                  return cId === childId;
                });
                return childComponent ? renderComponent(childComponent, newRenderedIds, componentId, containerPosition) : null;
              })
            ) : null}
          </ContainerComponent>
        );
      default:
        return null;
    }
  };

  if (loading) {
    return <div className="loading">Загрузка листа...</div>;
  }

  if (!sheet) {
    return <div className="error">Лист не найден</div>;
  }

  return (
    <div className="sheet-building-page">
      <div className="sheet-building-header">
        <div className="header-left">
          <button onClick={() => navigate('/sheets')} className="back-button">
            ← Назад
          </button>
          <h1>{sheet.Name}</h1>
        </div>
        <div className="header-actions">
          <div className="menu-container">
            <button onClick={() => setShowSheetMenu(!showSheetMenu)} className="menu-button">
              Меню листа
            </button>
            {showSheetMenu && (
              <div className="dropdown-menu">
                <button onClick={handleSaveSheet}>Сохранить</button>
                <button onClick={handleExportSheet}>Экспорт JSON</button>
                <button onClick={() => navigate(`/sheets`)}>Закрыть</button>
              </div>
            )}
          </div>
          <div className="menu-container">
            <button onClick={() => setShowComponentMenu(!showComponentMenu)} className="add-component-button">
              + Добавить компонент
            </button>
            {showComponentMenu && (
              <div className="dropdown-menu">
                <button onClick={() => handleCreateComponent(VarType.Number)}>Число</button>
                <button onClick={() => handleCreateComponent(VarType.Flag)}>Флаг</button>
                <button onClick={() => handleCreateComponent(VarType.Text)}>Текст</button>
                <button onClick={() => handleCreateComponent(VarType.Media)}>Медиа</button>
                <button onClick={() => handleCreateComponent(VarType.Container)}>Контейнер</button>
              </div>
            )}
          </div>
        </div>
      </div>

      {error && <div className="error-banner">{error}</div>}

      <div className="sheet-canvas" onClick={() => setSelectedComponents(new Set())}>
        {components
          .filter(c => {
            const sheetId = typeof sheet.Id === 'object' ? sheet.Id.Value : sheet.Id;
            const cSid = typeof c.SId === 'object' ? c.SId.Value : c.SId;
            const isChildOfContainer = components.some(container => {
              if (container.VarType !== VarType.Container) return false;
              try {
                if (container.VarValue && typeof container.VarValue === 'string') {
                  const parsed = JSON.parse(container.VarValue);
                  if (Array.isArray(parsed)) {
                    const containerId = typeof container.Id === 'object' ? container.Id.Value : container.Id;
                    const componentId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
                    return parsed.includes(componentId) && containerId !== componentId;
                  }
                }
              } catch {
                return false;
              }
              return false;
            });
            return cSid === sheetId && !isChildOfContainer;
          })
          .map(component => renderComponent(component))}
      </div>
    </div>
  );
};

