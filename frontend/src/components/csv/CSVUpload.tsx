import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import Papa from 'papaparse';
import { useAppState } from '../../hooks/useAppState';
// import { uploadFile } from '../../api/files'; // TODO: Implement API calls

export const CSVUpload: React.FC = () => {
  const { state, dispatch } = useAppState();

  const onDrop = useCallback((acceptedFiles: File[]) => {
    // 1. Check file limits
    if (state.csvFiles.length + acceptedFiles.length > 10) {
      dispatch({ type: 'SET_ERROR', payload: 'Maximum 10 files allowed.' });
      return;
    }

    acceptedFiles.forEach((file) => {
      // 2. Parse locally for preview (optional, or wait for backend)
      Papa.parse(file, {
        header: true,
        skipEmptyLines: true,
        complete: (results) => {
          console.log('Parsed:', results);
          // In a real app, we would upload here:
          // uploadFile(file, state.currentMonth).then(serverFile => ...);
          
          // For now, mock adding to state
          dispatch({
            type: 'ADD_FILE',
            payload: {
              id: Date.now(), // Temporary ID
              monthlyBudgetId: 0,
              filename: file.name,
              fileSize: file.size,
              uploadDate: new Date().toISOString(),
              rowCount: results.data.length,
            },
          });
        },
        error: (error) => {
          dispatch({ type: 'SET_ERROR', payload: `Failed to parse ${file.name}: ${error.message}` });
        }
      });
    });
  }, [state.csvFiles, dispatch]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'text/csv': ['.csv'],
      'application/vnd.ms-excel': ['.csv'],
    },
    maxSize: 10 * 1024 * 1024, // 10MB
  });

  return (
    <div
      {...getRootProps()}
      className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors
        ${isDragActive ? 'border-indigo-500 bg-indigo-50' : 'border-gray-300 hover:border-indigo-400'}
      `}
    >
      <input {...getInputProps()} />
      {isDragActive ? (
        <p className="text-indigo-600">Drop the CSV files here...</p>
      ) : (
        <div className="space-y-2">
          <p className="text-gray-600">Drag & drop CSV files here, or click to select</p>
          <p className="text-xs text-gray-400">Max 10 files, 10MB each</p>
        </div>
      )}
    </div>
  );
};

