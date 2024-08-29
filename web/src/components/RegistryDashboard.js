import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts';
import { AlertCircle, CheckCircle, PlusCircle } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

const RegistryDashboard = () => {
  const [registries, setRegistries] = useState([]);
  const [selectedRegistry, setSelectedRegistry] = useState(null);
  const [items, setItems] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchRegistries();
  }, []);

  const fetchRegistries = async () => {
    try {
      const response = await fetch('/registries');
      const data = await response.json();
      setRegistries(data);
    } catch (err) {
      setError('Failed to fetch registries');
    }
  };

  const fetchItems = async (registry) => {
    try {
      const response = await fetch(`/registry/${registry}/list`);
      const data = await response.json();
      setItems(data);
      setSelectedRegistry(registry);
    } catch (err) {
      setError(`Failed to fetch items for ${registry}`);
    }
  };

  const renderChart = () => {
    const data = registries.map(r => ({ name: r, count: items.length }));
    return (
      <LineChart width={600} height={300} data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="name" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="count" stroke="#8884d8" />
      </LineChart>
    );
  };

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">Registry Service Dashboard</h1>
      
      {error && (
        <Alert variant="destructive" className="mb-4">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div>
          <h2 className="text-xl font-semibold mb-2">Registries</h2>
          <ul className="space-y-2">
            {registries.map(registry => (
              <li 
                key={registry}
                className="flex items-center justify-between p-2 bg-gray-100 rounded cursor-pointer hover:bg-gray-200"
                onClick={() => fetchItems(registry)}
              >
                <span>{registry}</span>
                <CheckCircle className="h-4 w-4 text-green-500" />
              </li>
            ))}
          </ul>
          <button 
            className="mt-4 flex items-center px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            onClick={() => {/* Implement add registry functionality */}}
          >
            <PlusCircle className="h-4 w-4 mr-2" />
            Add Registry
          </button>
        </div>
        <div>
          {selectedRegistry && (
            <>
              <h2 className="text-xl font-semibold mb-2">Items in {selectedRegistry}</h2>
              <ul className="space-y-2">
                {items.map(item => (
                  <li key={item.id} className="p-2 bg-gray-100 rounded">
                    {item.id} - {item.type}
                  </li>
                ))}
              </ul>
            </>
          )}
        </div>
      </div>
      
      <div className="mt-8">
        <h2 className="text-xl font-semibold mb-2">Registry Overview</h2>
        {renderChart()}
      </div>
    </div>
  );
};

export default RegistryDashboard;