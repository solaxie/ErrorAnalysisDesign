document.addEventListener('DOMContentLoaded', () => {
    loadAttitudeButtons();
});

function loadAttitudeButtons() {
    fetch('/api/attitudes')
        .then(response => response.json())
        .then(attitudes => {
            const container = document.getElementById('attitude-buttons');
            attitudes.forEach(attitude => {
                const button = document.createElement('button');
                button.className = 'attitude-button';
                button.textContent = `${attitude.name} (${attitude.progress}/${attitude.total})`;
                button.setAttribute('hx-get', `/attitude/${attitude.name}`);
                button.setAttribute('hx-target', 'main');
                button.setAttribute('hx-swap', 'innerHTML');
                container.appendChild(button);
            });
        })
        .catch(error => console.error('Error loading attitudes:', error));
}

function setupImageSwipe(imageElement) {
    const hammer = new Hammer(imageElement);
    hammer.get('swipe').set({ direction: Hammer.DIRECTION_ALL });

    hammer.on('swipe', (event) => {
        const direction = event.direction;
        let action;

        switch (direction) {
            case Hammer.DIRECTION_RIGHT:
                action = 'correct';
                break;
            case Hammer.DIRECTION_LEFT:
                action = 'wrong';
                break;
            case Hammer.DIRECTION_DOWN:
                action = 'undo';
                break;
            case Hammer.DIRECTION_UP:
                action = 'save_exit';
                break;
            default:
                return;
        }

        const imageName = imageElement.getAttribute('data-image-name');
        const attitude = imageElement.getAttribute('data-attitude');

        fetch('/api/feedback', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ imageName, attitude, action }),
        })
        .then(response => response.json())
        .then(data => {
            if (action === 'save_exit') {
                window.location.href = '/';
            } else {
                // Load next image
                loadNextImage(attitude);
            }
        })
        .catch(error => console.error('Error sending feedback:', error));
    });
}

function loadNextImage(attitude) {
    fetch(`/api/next-image/${attitude}`)
        .then(response => response.json())
        .then(data => {
            const imageContainer = document.querySelector('.image-container');
            const img = document.createElement('img');
            img.src = data.imagePath;
            img.setAttribute('data-image-name', data.imageName);
            img.setAttribute('data-attitude', attitude);
            imageContainer.innerHTML = '';
            imageContainer.appendChild(img);
            setupImageSwipe(img);

            // Update attitude value
            document.querySelector('.attitude-value').textContent = data.attitudeValue;
        })
        .catch(error => console.error('Error loading next image:', error));
}

// Setup keyboard events for desktop version
document.addEventListener('keydown', (event) => {
    const img = document.querySelector('.image-container img');
    if (!img) return;

    let action;
    switch (event.key) {
        case 'ArrowRight':
            action = 'correct';
            break;
        case 'ArrowLeft':
            action = 'wrong';
            break;
        case 'ArrowDown':
            action = 'undo';
            break;
        case 'ArrowUp':
            action = 'save_exit';
            break;
        default:
            return;
    }

    const imageName = img.getAttribute('data-image-name');
    const attitude = img.getAttribute('data-attitude');

    fetch('/api/feedback', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ imageName, attitude, action }),
    })
    .then(response => response.json())
    .then(data => {
        if (action === 'save_exit') {
            window.location.href = '/';
        } else {
            // Load next image
            loadNextImage(attitude);
        }
    })
    .catch(error => console.error('Error sending feedback:', error));
});
