import { useContext } from 'react';
import { 
  Box, 
  Grid, 
  Paper, 
  Typography, 
  Card, 
  CardContent, 
  CardHeader,
  Divider,
  List,
  ListItem,
  ListItemText,
  Button,
} from '@mui/material';
import { AuthContext } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

const DashboardPage = () => {
  const { user } = useContext(AuthContext);
  const navigate = useNavigate();

  // This would be replaced with real data from API calls in a production app
  const dashboardData = {
    pendingOrders: 5,
    inventoryAlerts: 3,
    totalRecipes: 42,
    activeIngredients: 28,
    popularCocktails: [
      { id: 1, name: 'Mojito', orders: 128 },
      { id: 2, name: 'Old Fashioned', orders: 97 },
      { id: 3, name: 'Margarita', orders: 84 },
      { id: 4, name: 'Moscow Mule', orders: 76 },
      { id: 5, name: 'Whiskey Sour', orders: 69 },
    ],
    recentOrders: [
      { id: 1001, time: '2 minutes ago', items: 3, status: 'pending' },
      { id: 1000, time: '15 minutes ago', items: 2, status: 'completed' },
      { id: 999, time: '32 minutes ago', items: 4, status: 'completed' },
    ],
    lowStockItems: [
      { id: 12, name: 'Vodka', current: '120ml', minimum: '200ml' },
      { id: 8, name: 'Lime Juice', current: '50ml', minimum: '100ml' },
      { id: 5, name: 'Whiskey', current: '180ml', minimum: '200ml' },
    ]
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 4 }}>
        Welcome back, {user?.username}!
      </Typography>

      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Paper
            elevation={0}
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              bgcolor: 'primary.main',
              color: 'primary.contrastText',
              borderRadius: 2,
            }}
          >
            <Typography variant="h4">{dashboardData.pendingOrders}</Typography>
            <Typography variant="body2">Pending Orders</Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Paper
            elevation={0}
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              bgcolor: 'error.main',
              color: 'error.contrastText',
              borderRadius: 2,
            }}
          >
            <Typography variant="h4">{dashboardData.inventoryAlerts}</Typography>
            <Typography variant="body2">Inventory Alerts</Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Paper
            elevation={0}
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              bgcolor: 'secondary.main',
              color: 'secondary.contrastText',
              borderRadius: 2,
            }}
          >
            <Typography variant="h4">{dashboardData.totalRecipes}</Typography>
            <Typography variant="body2">Total Recipes</Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Paper
            elevation={0}
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              bgcolor: 'info.main',
              color: 'info.contrastText',
              borderRadius: 2,
            }}
          >
            <Typography variant="h4">{dashboardData.activeIngredients}</Typography>
            <Typography variant="body2">Active Ingredients</Typography>
          </Paper>
        </Grid>
      </Grid>

      {/* Detail Cards */}
      <Grid container spacing={3}>
        {/* Popular Cocktails */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardHeader title="Popular Cocktails" />
            <Divider />
            <CardContent>
              <List>
                {dashboardData.popularCocktails.map((cocktail) => (
                  <ListItem key={cocktail.id} disablePadding sx={{ py: 1 }}>
                    <ListItemText
                      primary={cocktail.name}
                      secondary={`${cocktail.orders} orders`}
                    />
                  </ListItem>
                ))}
              </List>
              <Button 
                variant="outlined" 
                fullWidth 
                sx={{ mt: 2 }}
                onClick={() => navigate('/recipes')}
              >
                View All Recipes
              </Button>
            </CardContent>
          </Card>
        </Grid>

        {/* Recent Orders */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardHeader title="Recent Orders" />
            <Divider />
            <CardContent>
              <List>
                {dashboardData.recentOrders.map((order) => (
                  <ListItem key={order.id} disablePadding sx={{ py: 1 }}>
                    <ListItemText
                      primary={`Order #${order.id}`}
                      secondary={`${order.time} • ${order.items} items • ${order.status}`}
                    />
                  </ListItem>
                ))}
              </List>
              <Button 
                variant="outlined" 
                fullWidth 
                sx={{ mt: 2 }}
                onClick={() => navigate('/orders')}
              >
                View All Orders
              </Button>
            </CardContent>
          </Card>
        </Grid>

        {/* Low Stock Items */}
        <Grid item xs={12} md={4}>
          <Card>
            <CardHeader title="Low Stock Items" />
            <Divider />
            <CardContent>
              <List>
                {dashboardData.lowStockItems.map((item) => (
                  <ListItem key={item.id} disablePadding sx={{ py: 1 }}>
                    <ListItemText
                      primary={item.name}
                      secondary={`Current: ${item.current} • Minimum: ${item.minimum}`}
                    />
                  </ListItem>
                ))}
              </List>
              <Button 
                variant="outlined" 
                fullWidth 
                sx={{ mt: 2 }}
                onClick={() => navigate('/ingredients')}
              >
                View Inventory
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default DashboardPage; 