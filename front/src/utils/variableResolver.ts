import { Component } from '../types';

export interface VariableContext {
  components: Component[];
  sheetId: number;
}

/**
 * Resolves variable references in a string (e.g., "$VarName", "$VarName[0]")
 * Returns the resolved value or the original string if no variables found
 */
export function resolveVariables(
  text: string,
  context: VariableContext,
  currentComponentId?: number
): string {
  const variablePattern = /\$(\w+)(?:\[(\d+)\])?/g;
  
  return text.replace(variablePattern, (match, varName, index) => {
    const component = context.components.find(c => {
      const cSid = typeof c.SId === 'object' ? c.SId.Value : c.SId;
      const cId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
      return c.VarName === varName && 
             cSid === context.sheetId &&
             cId !== currentComponentId;
    });
    
    if (!component) {
      return match;
    }

    try {
      if (typeof component.VarValue === 'string') {
        const parsed = JSON.parse(component.VarValue);
        
        if (Array.isArray(parsed)) {
          const idx = index !== undefined ? parseInt(index) : 0;
          if (idx >= 0 && idx < parsed.length) {
            return String(parsed[idx]);
          }
          return match;
        } else if (parsed.values && Array.isArray(parsed.values)) {
          const idx = index !== undefined ? parseInt(index) : 0;
          if (idx >= 0 && idx < parsed.values.length) {
            return String(parsed.values[idx]);
          }
          return match;
        } else {
          return String(parsed);
        }
      } else {
        const decoder = new TextDecoder();
        const decoded = decoder.decode(component.VarValue as Uint8Array);
        try {
          const parsed = JSON.parse(decoded);
          if (Array.isArray(parsed)) {
            const idx = index !== undefined ? parseInt(index) : 0;
            return idx >= 0 && idx < parsed.length ? String(parsed[idx]) : match;
          }
          return String(parsed);
        } catch {
          return decoded;
        }
      }
    } catch {
      return typeof component.VarValue === 'string' 
        ? component.VarValue 
        : String(component.VarValue);
    }
  });
}

/**
 * Evaluates a formula expression with variable references
 * Supports basic arithmetic: +, -, *, /, parentheses
 */
export function evaluateFormula(
  formula: string,
  context: VariableContext,
  currentComponentId?: number
): number {
  try {
    let resolvedFormula = resolveVariables(formula, context, currentComponentId);
    
    const variablePattern = /\$(\w+)(?:\[(\d+)\])?/g;
    resolvedFormula = resolvedFormula.replace(variablePattern, (match, varName, index) => {
      const component = context.components.find(c => {
        const cSid = typeof c.SId === 'object' ? c.SId.Value : c.SId;
        const cId = typeof c.Id === 'object' ? c.Id.Value : c.Id;
        return c.VarName === varName && 
               cSid === context.sheetId &&
               cId !== currentComponentId;
      });
      
      if (!component) return '0';
      
      try {
        if (typeof component.VarValue === 'string') {
          const parsed = JSON.parse(component.VarValue);
          if (parsed.values && Array.isArray(parsed.values)) {
            const idx = index !== undefined ? parseInt(index) : 0;
            return idx >= 0 && idx < parsed.values.length 
              ? String(parsed.values[idx]) 
              : '0';
          }
          return String(parsed);
        }
      } catch {
        return '0';
      }
      return '0';
    });
    
    const result = Function(`"use strict"; return (${resolvedFormula})`)();
    return typeof result === 'number' ? result : 0;
  } catch {
    return 0;
  }
}

/**
 * Checks if a string contains variable references
 */
export function hasVariables(text: string): boolean {
  return /\$\w+/.test(text);
}

/**
 * Gets all variable names referenced in a string
 */
export function getReferencedVariables(text: string): string[] {
  const variablePattern = /\$(\w+)/g;
  const variables = new Set<string>();
  let match;
  
  while ((match = variablePattern.exec(text)) !== null) {
    variables.add(match[1]);
  }
  
  return Array.from(variables);
}

