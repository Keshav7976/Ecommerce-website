// app.js - full updated version
const App = (function() {
  const API = 'http://localhost:8080';
  const tokenKey = 'jwt_token';

  // --- TOKEN ---
  function setToken(t) { localStorage.setItem(tokenKey, t); }
  function getToken() { return localStorage.getItem(tokenKey) || ''; }
  function clearToken() { localStorage.removeItem(tokenKey); }
  function authHeaders() {
    const t = getToken();
    return t ? { 'Authorization': 'Bearer ' + t } : {};
  }

  // --- NAVBAR ---
  function updateNavbar() {
    const token = getToken();
    const loginLink = document.getElementById('loginLink');
    const signupLink = document.getElementById('signupLink');
    const sellerLink = document.getElementById('sellerLink');
    const logoutLink = document.getElementById('logoutLink');

    if (loginLink) loginLink.style.display = token ? 'none' : 'inline-block';
    if (signupLink) signupLink.style.display = token ? 'none' : 'inline-block';
    if (sellerLink) sellerLink.style.display = token ? 'inline-block' : 'none';
    if (logoutLink) logoutLink.style.display = token ? 'inline-block' : 'none';

    if (logoutLink) {
      logoutLink.onclick = (e) => {
        e.preventDefault();
        clearToken();
        updateNavbar();
        window.location.href = 'index.html';
      };
    }
  }

  // --- ALERT ---
  function showAlert(elId, message, type = 'danger') {
    const el = document.getElementById(elId);
    if (!el) return;
    el.className = `alert alert-${type}`;
    el.innerText = message;
    el.classList.remove('d-none');
    setTimeout(() => el.classList.add('d-none'), 4000);
  }

  // --- HOME ---
  async function initHome() {
    updateNavbar();
    await loadCategoriesInto('#categoryFilter', true);
    bindHomeEvents();
    await loadItems();
  }

  function bindHomeEvents() {
    const filterBtn = document.getElementById('filterBtn');
    const clearBtn = document.getElementById('clearFilterBtn');

    filterBtn?.addEventListener('click', loadItems);
    clearBtn?.addEventListener('click', () => {
      document.getElementById('categoryFilter').value = '';
      document.getElementById('minPrice').value = '';
      document.getElementById('maxPrice').value = '';
      loadItems();
    });
  }

  async function loadItems() {
    const category = document.getElementById('categoryFilter')?.value || '';
    const min = document.getElementById('minPrice')?.value || '';
    const max = document.getElementById('maxPrice')?.value || '';

    const params = new URLSearchParams();
    if (category) params.append('category_id', category);
    if (min) params.append('minPrice', min);
    if (max) params.append('maxPrice', max);

    const row = document.getElementById('itemsRow');
    row.innerHTML = `<div class="text-center my-5 w-100">Loading...</div>`;

    try {
      const res = await fetch(`${API}/items?${params.toString()}`);
      if (!res.ok) throw new Error('Failed to fetch items');
      const items = await res.json();

      if (!Array.isArray(items) || items.length === 0) {
        row.innerHTML = `<div class="text-center my-5 w-100 text-muted">No items</div>`;
        return;
      }

      row.innerHTML = '';
      items.forEach(it => {
        const id = it.id ?? it.ID;
        const name = it.name ?? it.Name ?? 'Unnamed';
        const price = it.price ?? it.Price ?? 0;
        const image = it.image_url ?? it.ImageURL ?? 'https://placehold.co/400x300?text=No+Image';
        const categoryName = it.category ?? it.Category ?? '';

        const col = document.createElement('div');
        col.className = 'col-sm-6 col-md-4 col-lg-3';
        col.innerHTML = `
          <div class="card item-card h-100">
            <img src="${image}" class="card-img-top" alt="${name}" style="height:180px;object-fit:cover;">
            <div class="card-body d-flex flex-column">
              <h5 class="card-title">${escapeHtml(name)}</h5>
              <p class="text-muted mb-1">${escapeHtml(categoryName)}</p>
              <p class="price mb-2">₹${Number(price).toLocaleString()}</p>
              <div class="mt-auto d-grid">
                <button class="btn btn-success btn-sm add-to-cart" data-id="${id}">Add to Cart</button>
              </div>
            </div>
          </div>`;
        row.appendChild(col);
      });

      document.querySelectorAll('.add-to-cart').forEach(btn => {
        btn.addEventListener('click', async (e) => {
          const id = e.currentTarget.dataset.id;
          await addToCart(id, 1);
        });
      });

    } catch (err) {
      row.innerHTML = `<div class="text-center my-5 w-100 text-danger">Failed to load items</div>`;
      console.error(err);
    }
  }

  // --- CATEGORIES ---
  async function loadCategoriesInto(selector, includeBlank = false) {
    try {
      const res = await fetch(`${API}/categories`);
      if (!res.ok) return;
      const cats = await res.json();
      const sel = document.querySelector(selector);
      if (!sel) return;
      sel.innerHTML = includeBlank ? `<option value="">All categories</option>` : '';
      cats.forEach(c => {
        const id = c.ID ?? c.id;
        const name = c.name ?? c.Name ?? '';
        const opt = document.createElement('option');
        opt.value = id ?? name;
        opt.text = name;
        sel.appendChild(opt);
      });
    } catch (err) {
      console.warn('Could not load categories', err);
    }
  }

  // --- LOGIN ---
  async function initLogin() {
    updateNavbar();
    document.getElementById('loginBtn').addEventListener('click', async () => {
      const username = document.getElementById('loginUsername').value.trim();
      const password = document.getElementById('loginPassword').value;
      if (!username || !password) { showAlert('loginAlert','Provide username and password','warning'); return; }

      try {
        const res = await fetch(`${API}/login`, {
          method: 'POST',
          headers: {'Content-Type':'application/json'},
          body: JSON.stringify({ username, password })
        });
        const data = await res.json().catch(()=>({}));
        if (res.ok && data.token) {
          setToken(data.token);
          updateNavbar();
          showAlert('loginAlert','Login successful','success');
          setTimeout(()=> { window.location.href = 'index.html'; }, 1000);
        } else {
          showAlert('loginAlert', data.error || data.message || 'Login failed', 'danger');
        }
      } catch (err) {
        showAlert('loginAlert','Network error','danger');
      }
    });
  }

  // --- SIGNUP ---
  async function initSignup() {
    updateNavbar();
    document.getElementById('signupBtn').addEventListener('click', async () => {
      const username = document.getElementById('signupUsername').value.trim();
      const password = document.getElementById('signupPassword').value;
      if (!username || !password) { showAlert('signupAlert','Provide username and password','warning'); return; }

      try {
        const res = await fetch(`${API}/signup`, {
          method: 'POST',
          headers: {'Content-Type':'application/json'},
          body: JSON.stringify({ username, password })
        });
        const data = await res.json().catch(()=>({}));
        if (res.ok) {
          showAlert('signupAlert','Account created. Please login.','success');
          setTimeout(()=> { window.location.href = 'login.html'; }, 1000);
        } else {
          showAlert('signupAlert', data.error || data.message || 'Signup failed','danger');
        }
      } catch (err) {
        showAlert('signupAlert','Network error','danger');
      }
    });
  }

  // --- CART ---
  async function addToCart(itemId, qty = 1) {
    const token = getToken();
    if (!token) { alert('Please login'); window.location.href='login.html'; return; }

    try {
      const res = await fetch(`${API}/cart/add`, {
        method: 'POST',
        headers: Object.assign({'Content-Type':'application/json'}, authHeaders()),
        body: JSON.stringify({ item_id: Number(itemId), qty: Number(qty) })
      });
      if (res.ok) {
        alert('Added to cart');
      } else {
        const data = await res.json().catch(()=>({}));
        alert(data.error || 'Failed to add');
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function initCart() {
    updateNavbar();
    await renderCart();
  }

  async function renderCart() {
    const token = getToken();
    if (!token) { window.location.href='login.html'; return; }

    const listNode = document.getElementById('cartList');
    const summaryNode = document.getElementById('cartSummary');
    listNode.innerHTML = 'Loading...';
    summaryNode.innerHTML = '';

    try {
      const res = await fetch(`${API}/cart`, { headers: authHeaders() });
      if (!res.ok) throw new Error('Failed to load cart');
      const cart = await res.json();
      if (!cart.items || cart.items.length === 0) {
        listNode.innerHTML = '<div class="text-muted">Cart is empty</div>';
        return;
      }

      let total = 0;
      listNode.innerHTML = '';
      cart.items.forEach(it => {
        const subtotal = it.qty * it.price;
        total += subtotal;
        const li = document.createElement('div');
        li.className = 'list-group-item d-flex justify-content-between align-items-center';
        li.innerHTML = `
          <div>
            <strong>${escapeHtml(it.name)}</strong> x ${it.qty}
            <div class="small text-muted">₹${it.price.toLocaleString()} each</div>
          </div>
          <div>₹${subtotal.toLocaleString()}</div>`;
        listNode.appendChild(li);
      });
      summaryNode.innerHTML = `<strong>Total: ₹${total.toLocaleString()}</strong>`;
    } catch (err) {
      listNode.innerHTML = '<div class="text-danger">Failed to load cart</div>';
      console.error(err);
    }
  }

  // --- SELLER ---
  async function initSeller() {
    updateNavbar();
    const token = getToken();
    if (!token) { window.location.href='login.html'; return; }

    await loadCategoriesInto('#prodCategory', false);
    await loadSellerProducts();

    document.getElementById('addProdBtn').addEventListener('click', async () => {
      const name = document.getElementById('prodName').value.trim();
      const price = Number(document.getElementById('prodPrice').value);
      const image = document.getElementById('prodImage').value.trim();
      const categoryName = document.querySelector('#prodCategory option:checked')?.text || '';

      if (!name || !price || !categoryName) { showAlert('sellerAlert','Fill all fields','warning'); return; }

      try {
        const res = await fetch(`${API}/items`, {
          method: 'POST',
          headers: Object.assign({'Content-Type':'application/json'}, authHeaders()),
          body: JSON.stringify({ name, price, image_url: image, category_name: categoryName })
        });
        const data = await res.json().catch(()=>({}));
        if (res.ok) {
          showAlert('sellerAlert','Product added','success');
          loadSellerProducts();
        } else {
          showAlert('sellerAlert', data.error || 'Failed to add','danger');
        }
      } catch (err) {
        showAlert('sellerAlert','Network error','danger');
      }
    });
  }

  async function loadSellerProducts() {
    const node = document.getElementById('sellerProducts');
    node.innerHTML = 'Loading...';
    try {
      const res = await fetch(`${API}/items`, { headers: authHeaders() });
      if (!res.ok) throw new Error('Failed to load products');
      const items = await res.json();
      node.innerHTML = '';
      items.forEach(it => {
        const name = it.name ?? it.Name ?? 'Item';
        const id = it.id ?? it.ID;
        const price = it.price ?? it.Price;
        const li = document.createElement('div');
        li.className = 'list-group-item d-flex justify-content-between align-items-center';
        li.innerHTML = `
          <div>
            <strong>${escapeHtml(name)}</strong>
            <div class="small text-muted">₹${Number(price).toLocaleString()}</div>
          </div>
          <div>
            <button class="btn btn-sm btn-outline-danger" data-id="${id}" onclick="App.deleteItem(event)">Delete</button>
          </div>`;
        node.appendChild(li);
      });
    } catch (err) {
      node.innerHTML = '<div class="text-danger">Failed to load products</div>';
      console.error(err);
    }
  }

  async function deleteItemEvent(e) {
    const id = e.currentTarget.dataset.id;
    if (!confirm('Delete this item?')) return;
    try {
      const res = await fetch(`${API}/items/${id}`, { method:'DELETE', headers: authHeaders() });
      if (res.ok) {
        loadSellerProducts();
      } else {
        const d = await res.json().catch(()=>({}));
        alert(d.error || 'Delete failed');
      }
    } catch (err) {
      alert('Network error');
    }
  }

  // --- UTIL ---
  function escapeHtml(s) {
    if (!s) return '';
    return s.replaceAll('&','&amp;').replaceAll('<','&lt;').replaceAll('>','&gt;');
  }

  return {
    initHome, initLogin, initSignup, initCart, initSeller,
    addToCart, deleteItem: deleteItemEvent
  };
})();
