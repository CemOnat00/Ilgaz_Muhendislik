document.addEventListener('DOMContentLoaded', () => {
    // Grid View Toggle Logic
    const gridViewBtn = document.getElementById('gridViewBtn');
    const listViewBtn = document.getElementById('listViewBtn');
    const productGrid = document.getElementById('productGrid');

    if (gridViewBtn && listViewBtn && productGrid) {
        gridViewBtn.addEventListener('click', () => {
            if (!gridViewBtn.classList.contains('active')) {
                gridViewBtn.classList.add('active');
                listViewBtn.classList.remove('active');
                
                // Reset to CSS default grid behavior
                productGrid.style.gridTemplateColumns = '';
                
                // Keep product card flex direction natural
                const cards = productGrid.querySelectorAll('.product-card');
                cards.forEach(card => card.style.flexDirection = 'column');
            }
        });

        listViewBtn.addEventListener('click', () => {
            if (!listViewBtn.classList.contains('active')) {
                listViewBtn.classList.add('active');
                gridViewBtn.classList.remove('active');
                
                // Force to single column
                productGrid.style.gridTemplateColumns = '1fr';
                
                // Very basic style modification for list view look, 
                // in reality a full CSS class modifier would be better (.list-view)
                const cards = productGrid.querySelectorAll('.product-card');
                if (window.innerWidth > 640) {
                    cards.forEach(card => card.style.flexDirection = 'row');
                }
            }
        });
        
        // Handle resize layout changes when in list view
        window.addEventListener('resize', () => {
            if (listViewBtn.classList.contains('active')) {
                const cards = productGrid.querySelectorAll('.product-card');
                if (window.innerWidth > 640) {
                    cards.forEach(card => card.style.flexDirection = 'row');
                } else {
                    cards.forEach(card => card.style.flexDirection = 'column');
                }
            }
        });
    }
});
