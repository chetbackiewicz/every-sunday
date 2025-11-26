import type { AppState, Action } from '../types';

export const initialState: AppState = {
  user: null,
  currentMonth: new Date().toISOString().slice(0, 7), // Current YYYY-MM
  budget: null,
  csvFiles: [],
  categories: [],
  transactions: [],
  cellReferences: [],
  isLoading: false,
  error: null,
  activeCategoryId: null,
  isLinkingMode: false,
};

export function appReducer(state: AppState, action: Action): AppState {
  switch (action.type) {
    case 'SET_USER':
      return { ...state, user: action.payload };
    case 'SET_MONTH':
      return { ...state, currentMonth: action.payload };
    case 'SET_BUDGET':
      return { ...state, budget: action.payload };
    case 'SET_FILES':
      return { ...state, csvFiles: action.payload };
    case 'ADD_FILE':
      return { ...state, csvFiles: [action.payload, ...state.csvFiles] };
    case 'REMOVE_FILE':
      return {
        ...state,
        csvFiles: state.csvFiles.filter((f) => f.id !== action.payload),
        transactions: state.transactions.filter((t) => t.csvFileId !== action.payload), // Cascade delete from state
      };
    case 'SET_CATEGORIES':
      return { ...state, categories: action.payload };
    case 'UPDATE_CATEGORY':
      // Complex logic to update tree node in place might be needed, or just reload.
      // For now, assume we reload or simple map if flat list.
      // Since categories are a tree, we might need a helper to walk and update.
      // Or if we keep a flat list + derived tree, it's easier.
      // The type says CategoryNode[] so it's a tree.
      // Let's assume payload replaces the tree or specific node handling is external.
      // Ideally, we refetch categories on update for simplicity in MVP.
      return state; 
    case 'SET_TRANSACTIONS':
      return { ...state, transactions: action.payload };
    case 'UPDATE_TRANSACTION':
      return {
        ...state,
        transactions: state.transactions.map((t) =>
          t.id === action.payload.id ? action.payload : t
        ),
      };
    case 'SET_CELL_REFERENCES':
      return { ...state, cellReferences: action.payload };
    case 'ADD_CELL_REFERENCE':
      return { ...state, cellReferences: [...state.cellReferences, action.payload] };
    case 'REMOVE_CELL_REFERENCE':
      return {
        ...state,
        cellReferences: state.cellReferences.filter((r) => r.id !== action.payload),
      };
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload };
    case 'SET_ERROR':
      return { ...state, error: action.payload };
    case 'SET_LINKING_MODE':
      return {
        ...state,
        isLinkingMode: action.payload.active,
        activeCategoryId: action.payload.categoryId,
      };
    default:
      return state;
  }
}

