import React from 'react';
import { useAppState } from '../../hooks/useAppState';
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table';
import type { Transaction } from '../../types';

const columnHelper = createColumnHelper<Transaction>();

export const TransactionTable: React.FC = () => {
  const { state, dispatch } = useAppState();

  // Cell click handler for linking
  const handleCellClick = (tx: Transaction) => {
     if (state.isLinkingMode && state.activeCategoryId) {
         // Toggle link
         const existingLink = state.cellReferences.find(
             r => r.transactionId === tx.id && r.categoryId === Number(state.activeCategoryId)
         );

         if (existingLink) {
             dispatch({ type: 'REMOVE_CELL_REFERENCE', payload: existingLink.id });
         } else {
             dispatch({ 
                 type: 'ADD_CELL_REFERENCE', 
                 payload: {
                     id: Date.now(), // Temp ID
                     categoryId: Number(state.activeCategoryId),
                     transactionId: tx.id,
                     amount: tx.amount,
                     createdAt: new Date().toISOString()
                 }
             });
         }
     }
  };

  const columns = [
      columnHelper.accessor('transactionDate', {
        header: 'Date',
        cell: (info) => new Date(info.getValue()).toLocaleDateString(),
      }),
      columnHelper.accessor('description', {
        header: 'Description',
      }),
      columnHelper.accessor('category', {
        header: 'Category',
      }),
      columnHelper.accessor('amount', {
        header: 'Amount',
        cell: (info) => {
            const tx = info.row.original;
            // Check if linked to ANY category
            const isLinked = state.cellReferences.some(r => r.transactionId === tx.id);
            // Check if linked to ACTIVE category
            const isLinkedToActive = state.activeCategoryId && state.cellReferences.some(
                r => r.transactionId === tx.id && r.categoryId === Number(state.activeCategoryId)
            );
            
            return (
                <div 
                   className={`
                       cursor-pointer px-2 py-1 rounded
                       ${isLinkedToActive ? 'bg-indigo-100 text-indigo-800 font-bold' : ''}
                       ${!isLinkedToActive && isLinked ? 'bg-green-50 text-green-800' : ''}
                       ${state.isLinkingMode && !isLinkedToActive ? 'hover:bg-gray-100' : ''}
                   `}
                   onClick={() => handleCellClick(tx)}
                >
                    {info.getValue().toFixed(2)}
                    {isLinkedToActive && <span className="ml-1">✓</span>}
                </div>
            );
        },
      }),
    ];

  const table = useReactTable({
    data: state.transactions,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  if (state.transactions.length === 0) {
    return <div className="text-center py-8 text-gray-500">No transactions yet. Upload a CSV to get started.</div>;
  }

  return (
    <div className="overflow-x-auto shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
      <table className="min-w-full divide-y divide-gray-300">
        <thead className="bg-gray-50">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th
                  key={header.id}
                  scope="col"
                  className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody className="divide-y divide-gray-200 bg-white">
          {table.getRowModel().rows.map((row) => (
            <tr key={row.id}>
              {row.getVisibleCells().map((cell) => (
                <td
                  key={cell.id}
                  className="whitespace-nowrap px-3 py-4 text-sm text-gray-500"
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
