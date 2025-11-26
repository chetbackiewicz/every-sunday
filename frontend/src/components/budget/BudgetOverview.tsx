import React from 'react';
import { MonthSelector } from '../layout/MonthSelector';
import { CSVUpload } from '../csv/CSVUpload';
import { TransactionTable } from '../csv/TransactionTable';
import { CategoryTree } from './CategoryTree';

export const BudgetOverview: React.FC = () => {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Budget Overview</h2>
        <MonthSelector />
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Budget sections */}
        <div className="bg-white overflow-hidden shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Categories</h3>
          <CategoryTree />
        </div>
        <div className="bg-white overflow-hidden shadow rounded-lg p-6 col-span-2">
          <div className="flex justify-between items-center mb-4">
             <h3 className="text-lg font-medium text-gray-900">Transactions</h3>
             <div className="w-64">
                {/* Small upload area or modal trigger */}
             </div>
          </div>
          <div className="mb-6">
             <CSVUpload />
          </div>
          <TransactionTable />
        </div>
      </div>
    </div>
  );
};

