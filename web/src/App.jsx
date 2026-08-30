import React, { useState, useEffect } from 'react';
import { register, login, logout, createNote, updateNote, deleteNote, getNotes } from './script';
import { authEvents, LOGOUT_EVENT } from './auth_event';

export default function App() {
  const [page, setPage] = useState('login');
  const [user, setUser] = useState(null);
  const [notes, setNotes] = useState([]);
  const [currentNote, setCurrentNote] = useState(null);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [authInput, setAuthInput] = useState({ username: '', password: '' });

  useEffect(() => {
    authEvents.addEventListener(LOGOUT_EVENT, setLoggedOutState);
    return () => authEvents.removeEventListener(LOGOUT_EVENT, setLoggedOutState);
  }, []);

  useEffect(() => {
    if (page === 'list') {
      fetchNotes();
    }
  }, [page]);

  const fetchNotes = async () => {
    const result = await getNotes();
    if (result.success) {
      setNotes(result.data.notes);
    }
  };

  const handleLogin = async (e) => {
    e.preventDefault();
    
    const result = await login(authInput.username, authInput.password);
    if (result.success) {
      setUser(authInput.username);
      setAuthInput({ username: '', password: '' })
      setPage('list');
      
      return 
    }
    alert(result.data);
  };

  const handleRegister = async (e) => {
    e.preventDefault();

    const result = await register(authInput.username, authInput.password);
    if (result.success) {
      setUser(authInput.username);
      setAuthInput({ username: '', password: '' })
      setPage('list');
      
      return
    }
    alert(result.data);
  };

  const handleLogout = async () => {
    const result = await logout();
    if (result.success) {
      setLoggedOutState();
    }
  };

  const handleSaveNote = async (e) => {
    e.preventDefault();
    if (currentNote) {
      const result = await updateNote(currentNote.id, title, content);
      if (result.success) {
        setNotes(notes.map(n => n.id === currentNote.id ? { ...n, title, content } : n));
      }
    } else {
      const result = await createNote(title, content);
      if (result.success){
        setNotes([...notes, { id: result.data.id, content: result.data.content, title: result.data.title }]);
      }
    }
    setPage('list');
    setTitle('');
    setContent('');
    setCurrentNote(null);
  };

  const handleEdit = (note) => {
    setCurrentNote(note);
    setTitle(note.title);
    setContent(note.content);
    setPage('edit');
  };

  const handleDelete = async (id) => {
    const result = await deleteNote(id);
    if (result.success) {
      setNotes(notes.filter(n => n.id !== id));
    }
  };

  const setLoggedOutState = () => {
    setUser(null);
    setAuthInput({ username: '', password: '' });
    setPage('login');
  };

  return (
    <div style={{ fontFamily: 'sans-serif', margin: 0, padding: 0 }}>
      {user && (
        <nav style={{ background: '#333', color: '#fff', padding: '10px 20px', display: 'flex', justifyContent: 'space-between' }}>
          <div>
            <button onClick={() => setPage('list')} style={{ background: 'none', border: 'none', color: '#fff', cursor: 'pointer', marginRight: 15 }}>Notes</button>
            <button onClick={() => { setCurrentNote(null); setTitle(''); setContent(''); setPage('create'); }} style={{ background: 'none', border: 'none', color: '#fff', cursor: 'pointer' }}>New Note</button>
          </div>
          <button onClick={handleLogout} style={{ background: '#d9534f', border: 'none', color: '#fff', padding: '5px 10px', cursor: 'pointer' }}>Logout</button>
        </nav>
      )}

      <div style={{ padding: 20 }}>
        {page === 'login' && (
          <div>
            <h2>Login</h2>
            <form onSubmit={handleLogin}>
              <div>
                <input type="text" placeholder="Username" 
                       value={authInput.username} 
                       onChange={e => setAuthInput({...authInput, username: e.target.value})} 
                       required />
              </div>
              <div style={{ marginTop: 10 }}>
                <input type="password" placeholder="Password" 
                       value={authInput.password} 
                       onChange={e => setAuthInput({...authInput, password: e.target.value})} 
                       required />
                </div>
              <button type="submit" style={{ marginTop: 10 }}>Login</button>
            </form>
            <p style={{ marginTop: 15 }}>Need an account? 
              <button onClick={() => setPage('register')} 
                      style={{ background: 'none', border: 'none', color: 'blue', cursor: 'pointer', padding: 0 }}>Register
              </button>
            </p>
          </div>
        )}

        {page === 'register' && (
          <div>
            <h2>Register</h2>
            <form onSubmit={handleRegister}>
              <div>
                <input type="text" placeholder="Choose Username" 
                       value={authInput.username} 
                       onChange={e => setAuthInput({...authInput, username: e.target.value})} 
                       required />
              </div>
              <div style={{ marginTop: 10 }}>
                <input type="password" placeholder="Choose Password" 
                       value={authInput.password} 
                       onChange={e => setAuthInput({...authInput, password: e.target.value})} 
                       required />
              </div>
              <button type="submit" style={{ marginTop: 10 }}>Register & Login</button>
            </form>
            <p style={{ marginTop: 15 }}>Have an account? 
              <button onClick={() => setPage('login')} 
                      style={{ background: 'none', border: 'none', color: 'blue', cursor: 'pointer', padding: 0 }}>Login
              </button>
            </p>
          </div>
        )}

        {page === 'list' && (
          <div>
            <h2>Your Notes</h2>
            {!notes || notes.length === 0 ? <p>No notes found.</p> : notes.map(note => (
              <div key={note.id} style={{ border: '1px solid #ccc', padding: 10, marginBottom: 10, borderRadius: 4 }}>
                <h3>{note.title}</h3>
                <p>{note.content}</p>
                <button onClick={() => handleEdit(note)} style={{ marginRight: 10 }}>Edit</button>
                <button onClick={() => handleDelete(note.id)} style={{ background: '#d9534f', color: '#fff', border: 'none', padding: '4px 8px' }}>Delete</button>
              </div>
            ))}
          </div>
        )}

        {(page === 'create' || page === 'edit') && (
          <div>
            <h2>{page === 'create' ? 'Create Note' : 'Edit Note'}</h2>
            <form onSubmit={handleSaveNote}>
              <div>
                <input type="text" 
                       placeholder="Title" 
                       value={title} 
                       onChange={e => setTitle(e.target.value)} required style={{ width: '100%', maxWidth: 400, padding: 8 }} />
              </div>
              <div style={{ marginTop: 10 }}>
                <textarea placeholder="Content" 
                          value={content} 
                          onChange={e => setContent(e.target.value)} 
                          required rows={5} 
                          style={{ width: '100%', maxWidth: 400, padding: 8 }} />
              </div>
              <button type="submit" style={{ marginTop: 10 }}>Save</button>
            </form>
          </div>
        )}
      </div>
    </div>
  );
}
