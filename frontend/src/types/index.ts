export interface User {
  id: number;
  email: string;
  createdAt: string;
}

export interface MonthlyBudget {
  id: number;
  userId: number;
  month: string; // YYYY-MM
  name: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface Category {
  id: number;
  monthlyBudgetId: number;
  name: string;
  projected: number;
  manualActual: number | null;
  createdAt: string;
  updatedAt: string;
  // Frontend specific for tree structure
  children?: Category[];
  parentId?: number; // Helper for updates
  depth?: number;
  calculatedActual?: number; // Derived
}

// Tree node structure for react-arborist might need specific fields, 
// but we can map or use this if compatible.
// React-arborist expects: id, children (optional), isOpen (optional)
export interface CategoryNode extends Category {
  children?: CategoryNode[];
}

export interface CSVFile {
  id: number;
  monthlyBudgetId: number;
  filename: string;
  uploadDate: string;
  fileSize: number;
  rowCount: number;
}

export interface Transaction {
  id: number;
  csvFileId: number;
  rowNumber: number;
  transactionDate: string;
  postDate: string;
  description: string;
  category: string;
  type: string;
  amount: number;
  note: string | null;
  createdAt: string;
  updatedAt: string;
  // Derived/Frontend state
  linkedCategories?: CellReference[];
}

export interface CellReference {
  id: number;
  categoryId: number;
  transactionId: number;
  amount: number;
  createdAt: string;
}

export interface AppState {
  user: User | null;
  currentMonth: string; // YYYY-MM
  budget: MonthlyBudget | null;
  csvFiles: CSVFile[];
  categories: CategoryNode[];
  transactions: Transaction[];
  cellReferences: CellReference[];
  
  // UI State
  isLoading: boolean;
  error: string | null;
  activeCategoryId: string | null; // For linking mode
  isLinkingMode: boolean;
}

export type Action =
  | { type: 'SET_USER'; payload: User | null }
  | { type: 'SET_MONTH'; payload: string }
  | { type: 'SET_BUDGET'; payload: MonthlyBudget | null }
  | { type: 'SET_FILES'; payload: CSVFile[] }
  | { type: 'ADD_FILE'; payload: CSVFile }
  | { type: 'REMOVE_FILE'; payload: number }
  | { type: 'SET_CATEGORIES'; payload: CategoryNode[] }
  | { type: 'UPDATE_CATEGORY'; payload: Category }
  | { type: 'SET_TRANSACTIONS'; payload: Transaction[] }
  | { type: 'UPDATE_TRANSACTION'; payload: Transaction }
  | { type: 'SET_CELL_REFERENCES'; payload: CellReference[] }
  | { type: 'ADD_CELL_REFERENCE'; payload: CellReference }
  | { type: 'REMOVE_CELL_REFERENCE'; payload: number }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_LINKING_MODE'; payload: { active: boolean; categoryId: string | null } };

