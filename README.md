# Swipe Mission

Swipe Mission is a web application that allows users to identify whether photos and their descriptions are correct.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/swipe-mission.git
   cd swipe-mission
   ```

2. Build and run the application:
   ```
   docker-compose up --build
   ```

3. Open your web browser and navigate to `http://localhost:8080`

## Usage

- On the home page, you'll see buttons for different attitudes (e.g., age, gender).
- Click on an attitude to start reviewing images.
- For each image:
  - Swipe right or press the right arrow key if the description is correct.
  - Swipe left or press the left arrow key if the description is wrong.
  - Swipe down or press the down arrow key to undo the last action.
  - Swipe up or press the up arrow key to save and exit.

## Project Structure

- `backend/`: Contains the Go backend code
- `frontend/`: Contains the HTML, CSS, and JavaScript for the frontend
- `db/`: Contains the ScyllaDB initialization script
- `test_image/`: Contains the test images
- `test_result/`: Contains the test result text files
- `Dockerfile`: Defines the Docker image for the application
- `docker-compose.yml`: Defines the services (app and database) for running the application

## Development

To make changes to the application:

1. Modify the code in the `backend/` or `frontend/` directories as needed.
2. Rebuild and run the application using `docker-compose up --build`.

## Testing

Currently, the application uses test images and results located in the `test_image/` and `test_result/` directories. To test with different images or results, replace the files in these directories before building the Docker image.

## License

This project is licensed under the MIT License.
