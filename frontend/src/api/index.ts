import axios from 'axios';
import { MonthlyBudget, Category, Transaction, CSVFile, CellReference } from '../types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Helper for auth headers (placeholder)
const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

export const api = {
  budget: {
    get: (month: string) => axios.get<MonthlyBudget>(`${API_URL}/budgets/${month}`, { headers: getAuthHeader() }),
    create: (month: string, name: string) => axios.post<MonthlyBudget>(`${API_URL}/budgets`, { month, name }, { headers: getAuthHeader() }),
  },
  categories: {
    list: (month: string) => axios.get<{categories: Category[]}>(`${API_URL}/budgets/${month}/categories`, { headers: getAuthHeader() }),
    create: (month: string, data: Partial<Category>) => axios.post<Category>(`${API_URL}/budgets/${month}/categories`, data, { headers: getAuthHeader() }),
    update: (id: number, data: Partial<Category>) => axios.put<Category>(`${API_URL}/categories/${id}`, data, { headers: getAuthHeader() }),
    delete: (id: number) => axios.delete(`${API_URL}/categories/${id}`, { headers: getAuthHeader() }),
  },
  files: {
    upload: (month: string, file: File) => {
       const formData = new FormData();
       formData.append('file', file);
       return axios.post<{file: CSVFile, transaction_count: number}>(`${API_URL}/budgets/${month}/files`, formData, { 
           headers: { ...getAuthHeader(), 'Content-Type': 'multipart/form-data' } 
       });
    },
    list: (month: string) => axios.get<{files: CSVFile[]}>(`${API_URL}/budgets/${month}/files`, { headers: getAuthHeader() }),
  },
  transactions: {
      list: (month: string) => axios.get<{transactions: Transaction[]}>(`${API_URL}/budgets/${month}/transactions`, { headers: getAuthHeader() }),
  },
  references: {
      create: (categoryId: number, transactionId: number) => axios.post<CellReference>(`${API_URL}/budgets/current/references`, { categoryId, transactionId }, { headers: getAuthHeader() }),
      delete: (id: number) => axios.delete(`${API_URL}/budgets/current/references/${id}`, { headers: getAuthHeader() }),
  }
};

