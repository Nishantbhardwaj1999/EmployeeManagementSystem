document.addEventListener('DOMContentLoaded', function () {
    const employeeList = document.getElementById('employee-list');
    const employeeForm = document.getElementById('employee-form');

    // Function to fetch and display all employees
    function fetchEmployees() {
        fetch('http://localhost:8080/employees')
            .then(response => response.json())
            .then(data => {
                employeeList.innerHTML = ''; // Clear the existing list
                if (data && Array.isArray(data)) { // Check if data is not null or undefined and is an array
                    employeeList.innerHTML = data.map(employee => `
                        <div>
                            <p><strong>Name:</strong> ${employee.name}</p>
                            <p><strong>Position:</strong> ${employee.position}</p>
                            <p><strong>Age:</strong> ${employee.age}</p>
                            <button class="delete-btn" data-id="${employee._id}">Delete</button>
                            <button class="update-btn" data-id="${employee._id}">Update</button>
                        </div>
                    `).join('');
                } else {
                    console.error('Invalid data received:', data);
                }
            })
            .catch(error => {
                console.error('Error fetching employee data:', error);
            });
    }

    // Fetch and display all employees on page load
    fetchEmployees();

    // Add new employee
    employeeForm.addEventListener('submit', async function (event) {
        event.preventDefault();
        const name = document.getElementById('name').value;
        const position = document.getElementById('position').value;
        const age = parseInt(document.getElementById('age').value);

        const newEmployee = {
            name: name,
            position: position,
            age: age
        };

        try {
            const response = await fetch('http://localhost:8080/employees', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(newEmployee),
            });
            if (response.ok) {
                fetchEmployees();
                employeeForm.reset();
            } else {
                throw new Error('Failed to add new employee');
            }
        } catch (error) {
            console.error('Error adding new employee:', error);
        }
    });

    // Update employee
    employeeList.addEventListener('click', async function (event) {
        if (event.target.classList.contains('update-btn')) {
            const employeeId = event.target.dataset.id;
            const updatedName = prompt('Enter updated name:');
            const updatedPosition = prompt('Enter updated position:');
            const updatedAge = parseInt(prompt('Enter updated age:'));

            const updatedEmployee = {
                name: updatedName,
                position: updatedPosition,
                age: updatedAge
            };

            try {
                const response = await fetch(`http://localhost:8080/employees/${employeeId}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(updatedEmployee),
                });
                if (response.ok) {
                    fetchEmployees();
                } else {
                    throw new Error('Failed to update employee');
                }
            } catch (error) {
                console.error('Error updating employee:', error);
            }
        }
    });

    // Delete employee
    // Delete employee
employeeList.addEventListener('click', async function (event) {
    if (event.target.classList.contains('delete-btn')) {
        const employeeId = event.target.dataset.id;
        try {
            const response = await fetch(`http://localhost:8080/employees/${employeeId}`, {
                method: 'DELETE',
            });
            if (response.ok) {
                fetchEmployees();
            } else {
                throw new Error('Failed to delete employee');
            }
        } catch (error) {
            console.error('Error deleting employee:', error);
        }
    }
});

});
